use actix_web::{get, post, web, App, HttpResponse, HttpServer, Responder};
use serde::{Deserialize, Serialize};

use reqwest::Client;
use std::time::Duration;


#[derive(Debug, Deserialize, Serialize)]
struct WeatherData {
    description: String,
    country: String,
    weather: String,
}

#[get("/")]
async fn index() -> impl Responder {
    HttpResponse::Ok().body("API de Rust en ejecución")
}



#[post("/input")]
async fn receive_weather(data: web::Json<WeatherData>) -> impl Responder {
    println!("Received: {:?}", data);
    // Enviar datos al servicio Go (API REST)
    match send_to_go_service(&data.into_inner()).await {
        Ok(_) => HttpResponse::Ok().body("Datos recibidos y enviados a Go"),
        Err(e) => {
            eprintln!("Error enviando a Go: {}", e);
            HttpResponse::InternalServerError().body("Error interno")
        }
    }
}

// Cliente HTTP para enviar datos a Go
async fn send_to_go_service(data: &WeatherData) -> Result<(), reqwest::Error> {
    let client = Client::new();
    let go_service_url = std::env::var("GO_SERVICE_URL")
        .unwrap_or_else(|_| "http://grpc-client:8081/".to_string());

    client
        .post(&go_service_url)
        .json(data)
        .timeout(Duration::from_secs(2))
        .send()
        .await?;

    Ok(())
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    println!("Starting server on http://0.0.0.0:8080");
    HttpServer::new(|| {
        App::new()
            .service(index)
            .service(receive_weather)
    })
    .bind(("0.0.0.0", 8080))?
    .run()
    .await
}
