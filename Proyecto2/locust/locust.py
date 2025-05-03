from locust import HttpUser, task, between
import random

class WeatherTrafficUser(HttpUser):
    wait_time = between(0.5, 2.5)  # Tiempo de espera entre peticiones

    # Descripciones y tipos de clima
    descriptions = [
        "Está lloviendo fuertemente",
        "Llovizna ligera",
        "Cielos despejados",
        "Nublado por completo",
        "Sol con nubes dispersas",
        "Tormenta eléctrica",
        "Niebla matutina",
        "Vientos fuertes",
        "Granizo ligero",
        "Cielo parcialmente nublado"
    ]

    # Lista de países comunes con sus abreviaciones
    countries = [
        "GT",  # Guatemala
        "MX",  # México
        "US",  # Estados Unidos
        "BR",  # Brasil
        "CO",  # Colombia
        "AR",  # Argentina
        "PE",  # Perú
        "CL",  # Chile
        "ES"  # España
    ]


    weather_types = ["Lluvioso", "Nubloso", "Soleado"]

    @task
    def send_weather_data(self):
        # Datos aleatorios para la petición
        payload = {
            "description": random.choice(self.descriptions),
            "country": random.choice(self.countries),
            "weather": random.choice(self.weather_types)
        }

        headers = {"Content-Type": "application/json"}

        # Enviar petición POST
        self.client.post(
            "/input",
            json=payload,
            headers=headers,
            name="Enviar datos climáticos"
        )