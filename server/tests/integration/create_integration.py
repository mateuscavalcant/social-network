import unittest
import requests

class SignupIntegratedTest(unittest.TestCase):

    def test_signup_success(self):
        url = 'http://localhost:8000/signup'
        data = {
            'name': 'John Doe',
            'email': 'john1@example.com',
            'password': 'password123',
            'confirm_password': 'password123'
        }
        response = requests.post(url, data=data)
        self.assertEqual(response.status_code, 200)


    def test_signup_invalid_data(self):
        url = 'http://localhost:8000/signup'
        data = {
            'name': '',
            'email': 'invalid_email',
            'password': 'short',
            'confirm_password': 'password'
        }
        response = requests.post(url, data=data)
        self.assertEqual(response.status_code, 400)


if __name__ == '__main__':
    unittest.main()