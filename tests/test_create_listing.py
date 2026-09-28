"""Test POST /listings endpoint."""
import requests
import time

TARGET_URL = "https://olx-api-58x8.onrender.com"


def test_create_listing_returns_201():
    payload = {
        "title": f"Test Listing {int(time.time())}",
        "description": "A test listing created by TestSprite",
        "price": 100,
        "city": "Test City"
    }
    r = requests.post(f"{TARGET_URL}/listings", json=payload)
    assert r.status_code == 201, f"Expected 201, got {r.status_code}"


def test_create_listing_returns_message():
    payload = {
        "title": f"Test Listing {int(time.time())}",
        "description": "A test listing created by TestSprite",
        "price": 100,
        "city": "Test City"
    }
    r = requests.post(f"{TARGET_URL}/listings", json=payload)
    data = r.json()
    assert data["message"] == "Listing created successfully", \
        f"Expected 'Listing created successfully', got '{data.get('message')}'"


def test_create_listing_returns_listing_object():
    payload = {
        "title": f"Test Listing {int(time.time())}",
        "description": "A test listing created by TestSprite",
        "price": 100,
        "city": "Test City"
    }
    r = requests.post(f"{TARGET_URL}/listings", json=payload)
    data = r.json()
    assert "listing" in data, "Response missing 'listing' object"
    listing = data["listing"]
    assert "id" in listing, "Listing missing 'id' field"
    assert "title" in listing, "Listing missing 'title' field"
    assert "created_at" in listing, "Listing missing 'created_at' field"


def test_create_listing_missing_title_returns_400():
    payload = {
        "description": "A test listing with no title",
        "price": 100,
        "city": "Test City"
    }
    r = requests.post(f"{TARGET_URL}/listings", json=payload)
    assert r.status_code == 400, f"Expected 400, got {r.status_code}"


def test_create_listing_invalid_price_returns_400():
    payload = {
        "title": "Test Listing",
        "description": "A test listing with invalid price",
        "price": -10,
        "city": "Test City"
    }
    r = requests.post(f"{TARGET_URL}/listings", json=payload)
    assert r.status_code == 400, f"Expected 400, got {r.status_code}"


def test_create_listing_missing_city_returns_400():
    payload = {
        "title": "Test Listing",
        "description": "A test listing with no city",
        "price": 100
    }
    r = requests.post(f"{TARGET_URL}/listings", json=payload)
    assert r.status_code == 400, f"Expected 400, got {r.status_code}"
