import unittest
import requests

class LoginPerformanceTest(unittest.TestCase):

    def test_login_success(self):
        url = 'http://localhost:8000/login'
        for i in range(0, 1000):
            # Gerar um endereço de email único
            email = f"test{i}@example.com"

            data = {
                'email': email,
                'password': 'password123'
            }
            response = requests.post(url, data=data)
            self.assertEqual(response.status_code, 200)

    def test_login_invalid(self):
        url = 'http://localhost:8000/login'
        for i in range(0, 1000):
            # Gerar um endereço de email único
            email = f"test{i}@example.com"

            data = {
                'email': email,
                'password': 'password'
            }
            response = requests.post(url, data=data)
            self.assertEqual(response.status_code, 401)


if __name__ == '__main__':
    unittest.main()
