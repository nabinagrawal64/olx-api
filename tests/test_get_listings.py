"""Test GET /listings endpoint."""
import requests

TARGET_URL = "https://olx-api-58x8.onrender.com"


def test_get_listings_returns_200():
    r = requests.get(f"{TARGET_URL}/listings")
    assert r.status_code == 200, f"Expected 200, got {r.status_code}"


def test_get_listings_returns_array():
    r = requests.get(f"{TARGET_URL}/listings")
    data = r.json()
    assert isinstance(data, list), f"Expected list, got {type(data).__name__}"


def test_get_listings_content_type_json():
    r = requests.get(f"{TARGET_URL}/listings")
    assert "application/json" in r.headers.get("Content-Type", ""), \
        f"Expected JSON content type, got '{r.headers.get('Content-Type')}'"


def test_get_listings_items_have_required_fields():
    r = requests.get(f"{TARGET_URL}/listings")
    data = r.json()
    if len(data) > 0:
        item = data[0]
        required_fields = ["id", "title", "description", "price", "city", "created_at"]
        for field in required_fields:
            assert field in item, f"Missing field '{field}' in listing item"
