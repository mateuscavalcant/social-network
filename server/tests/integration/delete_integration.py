import unittest
import requests

class DeleteIntegratedTest(unittest.TestCase):

    def test_delete_success(self):
        url = 'http://localhost:8000/delete-account'
        data = {
            'email': 'john1@example.com',
            'password': 'password123',
            'confirm_password': 'password123'
        }
        response = requests.delete(url, json=data)
        self.assertEqual(response.status_code, 200)


    def test_delete_invalid_data(self):
        url = 'http://localhost:8000/delete-account'
        data = {
            'email': 'invalid_email',
            'password': 'short',
            'confirm_password': 'password123'
        }
        response = requests.delete(url, json=data)
        self.assertEqual(response.status_code, 400)


if __name__ == '__main__':
    unittest.main()