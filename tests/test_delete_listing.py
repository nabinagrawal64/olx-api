"""Test DELETE /listings/{id} endpoint."""
import requests
import time

TARGET_URL = "https://olx-api-58x8.onrender.com"


def _create_listing():
    """Helper to create a listing and return its ID."""
    payload = {
        "title": f"Test Listing to Delete {int(time.time())}",
        "description": "A test listing to be deleted",
        "price": 50,
        "city": "Test City"
    }
    r = requests.post(f"{TARGET_URL}/listings", json=payload)
    assert r.status_code == 201, f"Failed to create listing: {r.status_code}"
    return r.json()["listing"]["id"]


def test_delete_listing_returns_200():
    listing_id = _create_listing()
    r = requests.delete(f"{TARGET_URL}/listings/{listing_id}")
    assert r.status_code == 200, f"Expected 200, got {r.status_code}"


def test_delete_listing_returns_message():
    listing_id = _create_listing()
    r = requests.delete(f"{TARGET_URL}/listings/{listing_id}")
    data = r.json()
    assert data["message"] == "Listing deleted successfully", \
        f"Expected 'Listing deleted successfully', got '{data.get('message')}'"


def test_delete_listing_removes_from_listings():
    listing_id = _create_listing()
    # Delete the listing
    requests.delete(f"{TARGET_URL}/listings/{listing_id}")
    # Verify it's gone from the listings
    r = requests.get(f"{TARGET_URL}/listings")
    data = r.json()
    ids = [item["id"] for item in data]
    assert listing_id not in ids, f"Listing {listing_id} still exists after deletion"
