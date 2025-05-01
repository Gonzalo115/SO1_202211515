use actix_web::{post, web, App, HttpResponse, HttpServer, Responder};
use serde::Deserialize;

#[derive(Debug, Deserialize)]
struct WeatherData {
    description: String,
    country: String,
    weather: String,
}

#[post("/input")]
async fn receive_weather(data: web::Json<WeatherData>) -> impl Responder {
    println!("Received: {:?}", data);
    // Proxima logica para enviar los datos a la API rest de go
    HttpResponse::Ok().body("Received weather data")
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    println!("Starting server on http://0.0.0.0:8080");
    HttpServer::new(|| {
        App::new()
            .service(receive_weather)
    })
    .bind(("0.0.0.0", 8080))?
    .run()
    .await
}
