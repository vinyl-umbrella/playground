import pytest


@pytest.mark.parametrize(
    "status_code, expected",
    [
        pytest.param(200, 200, id="200 OK"),
        pytest.param(404, 404, id="404 Not Found"),
        pytest.param(500, 500, id="500 Internal Server Error"),
    ],
)
def test_http_status_codes(status_code: int, expected: int):
    from samplepoetry.main import get_http_status

    res = get_http_status(status_code)
    print(res.status_code)
    assert res.status_code == expected
