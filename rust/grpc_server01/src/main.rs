use tonic::{transport::Server, Request, Response, Status};

mod user {
    tonic::include_proto!("example.user");
}

use user::{
    user_service_server::{UserService, UserServiceServer},
    Book, MyRequest, User,
};

#[derive(Default)]
pub struct MyUserService {}

#[tonic::async_trait]
impl UserService for MyUserService {
    async fn create_user(&self, request: Request<MyRequest>) -> Result<Response<User>, Status> {
        let b1 = Book {
            title: "book1".to_string(),
            author: "author1".to_string(),
        };
        let b2 = Book {
            title: "book2".to_string(),
            author: "author2".to_string(),
        };
        let books = vec![b1, b2];

        let req = request.into_inner();

        let reply = user::User {
            name: format!("{}", req.name).into(),
            age: req.age,
            books: books,
        };

        Ok(Response::new(reply))
    }
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "[::1]:8080".parse().unwrap();
    let user = MyUserService::default();

    Server::builder()
        .add_service(UserServiceServer::new(user))
        .serve(addr)
        .await?;
    Ok(())
}
