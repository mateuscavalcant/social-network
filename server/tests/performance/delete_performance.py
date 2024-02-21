import unittest
import requests

class DeletePerformanceTest(unittest.TestCase):
    def test_delete_success(self):
        url = 'http://localhost:8000/delete-account'
        for i in range(0, 1000):
            # Gerar um endereço de email único
            email = f"test{i}@example.com"

            data = {
                'email': email,
                'password': 'password123',
                'confirm_password': 'password123'
            }
            response = requests.delete(url, json=data)
            self.assertEqual(response.status_code, 200)

    def test_delete_invalid(self):
        url = 'http://localhost:8000/delete-account'
        for i in range(0, 1000):
            # Gerar um endereço de email único
            email = f"test{i}@example.com"

            data = {
                'email': email,
                'password': 'password',
                'confirm_password': 'password123'
            }
            response = requests.delete(url, json=data)
            self.assertEqual(response.status_code, 400)




if __name__ == '__main__':
    unittest.main()
