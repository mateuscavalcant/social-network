import unittest
import requests

class LoginIntegratedTest(unittest.TestCase):

    def test_login_success(self):
        url = 'http://localhost:8000/login'
        data = {
            'email': 'john1@example.com',
            'password': 'password123'
        }
        response = requests.post(url, data=data)
        self.assertEqual(response.status_code, 200)


    def test_login_invalid_data(self):
        url = 'http://localhost:8000/login'
        data = {
            'email': 'invalid_email',
            'password': 'short'
        }
        response = requests.post(url, data=data)
        self.assertEqual(response.status_code, 401)


if __name__ == '__main__':
    unittest.main()