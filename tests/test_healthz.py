"""Test GET /healthz endpoint."""
import requests

TARGET_URL = "https://olx-api-58x8.onrender.com"


def test_healthz_returns_200():
    r = requests.get(f"{TARGET_URL}/healthz")
    assert r.status_code == 200, f"Expected 200, got {r.status_code}"


def test_healthz_returns_ok_status():
    r = requests.get(f"{TARGET_URL}/healthz")
    data = r.json()
    assert data["status"] == "okay", f"Expected 'okay', got '{data.get('status')}'"


def test_healthz_content_type_json():
    r = requests.get(f"{TARGET_URL}/healthz")
    assert "application/json" in r.headers.get("Content-Type", ""), \
        f"Expected JSON content type, got '{r.headers.get('Content-Type')}'"
