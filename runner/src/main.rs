#![feature(custom_derive, plugin)]
#![plugin(serde_macros)]

#[macro_use] extern crate mrusty;
extern crate rmp;
extern crate rmp_serde;
extern crate rmp_rpc;
extern crate byteorder;
extern crate serde;

mod runner;
mod protocol;
mod api;

use std::os::unix::io::FromRawFd;
use std::os::unix::net::UnixStream;
use std::fs::File;
use byteorder::{NativeEndian, ReadBytesExt};

use runner::Runner;

pub fn main() {
    let mut connection = unsafe {
        UnixStream::from_raw_fd(3)
    };

    let seed = {
        let mut random = unsafe { File::from_raw_fd(4) };
        random.read_i32::<NativeEndian>().unwrap()
    };

    let runner = Runner::new(seed).unwrap();
    let mut server = rmp_rpc::Server::new(runner);

    server.serve(&mut connection);
}
