use nix::unistd::{fork, getpid, ForkResult};

fn main() {
    println!("fork test");

    match unsafe { fork() } {
        Ok(ForkResult::Parent { child, .. }) => {
            println!("## parent pid={}: child is {}", getpid(), child)
        }
        Ok(ForkResult::Child) => {
            println!("## child pid={}", getpid())
        }
        Err(err) => println!("error {:?}", err),
    }
}
