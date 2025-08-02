#include <aws/core/Aws.h>
#include <aws/core/auth/AWSCredentialsProvider.h>
#include <aws/core/utils/stream/ResponseStream.h>
#include <aws/core/utils/stream/PreallocatedStreamBuf.h>
#include <aws/core/utils/threading/Executor.h>
#include <aws/core/utils/memory/AWSMemory.h>
#include <aws/s3/S3Client.h>
#include <aws/s3/model/GetObjectRequest.h>
#include <aws/s3/model/ListObjectsV2Request.h>
#include <aws/s3/model/Object.h>
#include <aws/transfer/TransferManager.h>
#include <aws/transfer/TransferHandle.h>
#include <iostream>
#include <fstream>
#include <memory>
#include <string>
#include <unistd.h>
#include <fcntl.h>
#include <cassert>

using namespace std;

// AWS公式サンプルに基づいたメモリストリーム実装
class MyUnderlyingStream : public Aws::IOStream
{
public:
    using Base = Aws::IOStream;
    // Provide a customer-controlled streambuf to hold data from the bucket.
    explicit MyUnderlyingStream(std::streambuf *buf)
        : Base(buf)
    {
    }

    ~MyUnderlyingStream() override = default;
};

class S3Stream
{
private:
    shared_ptr<Aws::S3::S3Client> s3Client_;
    shared_ptr<Aws::Transfer::TransferManager> transferManager_;

public:
    S3Stream(const string &region)
    {
        // AWS設定の初期化
        Aws::SDKOptions options;
        Aws::InitAPI(options);

        Aws::Client::ClientConfiguration config;
        config.region = region;

        cout << "Using S3 region: " << region << endl;
        s3Client_ = make_shared<Aws::S3::S3Client>(config);

        // Transfer Manager の設定
        auto executor = Aws::MakeShared<Aws::Utils::Threading::PooledThreadExecutor>("executor", 25);
        Aws::Transfer::TransferManagerConfiguration transfer_config(executor.get());
        transfer_config.s3Client = s3Client_;
        transfer_config.bufferSize = 2 * 1024 * 1024; // 2MB buffer

        transferManager_ = Aws::Transfer::TransferManager::Create(transfer_config);

        cout << "Transfer Manager initialized" << endl;
    }

    ~S3Stream()
    {
        Aws::SDKOptions options;
        Aws::ShutdownAPI(options);
    }

    // S3 List Objects
    bool listObjects(const string &bucketName, const string &prefix = "")
    {
        Aws::S3::Model::ListObjectsV2Request request;
        request.SetBucket(bucketName);
        if (!prefix.empty())
        {
            request.SetPrefix(prefix);
        }

        auto outcome = s3Client_->ListObjectsV2(request);

        if (!outcome.IsSuccess())
        {
            const auto &error = outcome.GetError();
            cerr << "Error listing objects: " << error.GetMessage() << endl;
            cerr << "Error code: " << error.GetExceptionName() << endl;
            cerr << "HTTP response code: " << static_cast<int>(error.GetResponseCode()) << endl;

            if (error.GetExceptionName() == "PermanentRedirect")
            {
                cerr << "Hint: The bucket might be in a different region. " << endl;
                cerr << "Current region setting: ap-northeast-1" << endl;
                cerr << "Please verify the bucket location or update the region setting." << endl;
            }
            return false;
        }

        const auto &objects = outcome.GetResult().GetContents();
        cout << "Found " << objects.size() << " objects in bucket: " << bucketName << endl;

        for (const auto &object : objects)
        {
            cout << "  - " << object.GetKey() << " (Size: " << object.GetSize() << " bytes)" << endl;
        }

        return true;
    }

    bool downloadObject(const string &bucketName, const string &objectKey,
                        const string &localPath = "")
    {
        cerr << "Starting S3 Transfer Manager download..." << endl;

        const size_t BUFFER_SIZE = 10 * 1024 * 1024; // 10MiB
        Aws::Utils::Array<unsigned char> buffer(BUFFER_SIZE);

        // Create buffer to hold data received by the data stream.
        // The local variable 'streamBuffer' is captured by reference in a lambda.
        // It must persist until all downloading by the 'transfer_manager' is complete.
        Aws::Utils::Stream::PreallocatedStreamBuf streamBuffer(buffer.GetUnderlyingData(), buffer.GetLength());

        auto downloadHandle = transferManager_->DownloadFile(bucketName, objectKey,
                                                             [&streamBuffer]() -> Aws::IOStream *
                                                             { // Define a lambda expression for the callback method parameter to stream back the data.
                                                                 return Aws::New<MyUnderlyingStream>("TestTag", &streamBuffer);
                                                             });

        downloadHandle->WaitUntilFinished(); // Block calling thread until download is complete.
        auto downStat = downloadHandle->GetStatus();
        if (downStat != Aws::Transfer::TransferStatus::COMPLETED)
        {
            auto err = downloadHandle->GetLastError();
            cerr << "File download failed: " << err.GetMessage() << endl;
            return false;
        }
        cerr << "File download to memory finished." << endl;

        // Verify the download retrieved the expected length of data.
        assert(downloadHandle->GetBytesTotalSize() == downloadHandle->GetBytesTransferred());

        // Write the buffered data to output
        if (localPath.empty())
        {
            // 標準出力に書き込み
            ssize_t written = write(STDOUT_FILENO, buffer.GetUnderlyingData(),
                                    static_cast<size_t>(downloadHandle->GetBytesTransferred()));
            if (written != downloadHandle->GetBytesTransferred())
            {
                cerr << "Error writing to stdout" << endl;
                return false;
            }
        }
        else
        {
            // ファイルに書き込み
            Aws::OFStream storeFile(localPath.c_str(), Aws::OFStream::out | Aws::OFStream::trunc);
            if (!storeFile.is_open())
            {
                cerr << "Error opening output file: " << localPath << endl;
                return false;
            }
            storeFile.write((const char *)(buffer.GetUnderlyingData()),
                            static_cast<std::streamsize>(downloadHandle->GetBytesTransferred()));
            storeFile.close();
        }

        cerr << "File dump to output finished. Total bytes: " << downloadHandle->GetBytesTransferred() << endl;
        return true;
    }
};

int main(int argc, char *argv[])
{
    cout << "S3 Stream Processor starting..." << endl;

    if (argc < 3)
    {
        cout << "Usage: " << argv[0] << " <bucket-name> <command> [options]" << endl;
        cout << "Commands:" << endl;
        cout << "  list [prefix]                    - List objects in bucket" << endl;
        cout << "  get <object-key> [local-path]    - Download object (stream to console if no path)" << endl;
        return 1;
    }

    string bucketName = argv[1];
    string command = argv[2];

    try
    {
        // AWS認証情報の確認
        const char *accessKey = getenv("AWS_ACCESS_KEY_ID");
        const char *secretKey = getenv("AWS_SECRET_ACCESS_KEY");

        if (!accessKey || !secretKey)
        {
            cout << "Warning: AWS_ACCESS_KEY_ID or AWS_SECRET_ACCESS_KEY not found" << endl;
            cout << "Please set these environment variables for authentication" << endl;
        }

        // リージョンを環境変数から取得、デフォルトは ap-northeast-1
        string region = "ap-northeast-1";
        if (const char *envRegion = getenv("AWS_DEFAULT_REGION"))
        {
            region = envRegion;
        }
        else if (const char *envRegion = getenv("AWS_REGION"))
        {
            region = envRegion;
        }

        S3Stream processor(region);

        if (command == "list")
        {
            string prefix = argc > 3 ? argv[3] : "";
            processor.listObjects(bucketName, prefix);
        }
        else if (command == "get" && argc >= 4)
        {
            string objectKey = argv[3];
            string localPath = argc > 4 ? argv[4] : "";
            processor.downloadObject(bucketName, objectKey, localPath);
        }
        else
        {
            cerr << "Invalid command or missing arguments" << endl;
            return 1;
        }
    }
    catch (const exception &e)
    {
        cerr << "Exception: " << e.what() << endl;
        return 1;
    }

    return 0;
}
