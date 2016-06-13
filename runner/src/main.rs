#![feature(custom_derive, plugin)]
#![plugin(serde_macros)]

#[macro_use] extern crate mrusty;
extern crate rmp;
extern crate rmp_serde;
extern crate rmp_rpc;
extern crate byteorder;
extern crate serde;
extern crate getopts;

mod runner;
mod protocol;
mod api;

use std::os::unix::io::FromRawFd;
use std::os::unix::net::UnixStream;
use std::fs::File;
use byteorder::{NativeEndian, ReadBytesExt};
use getopts::Options;

use runner::Runner;

#[cfg(feature = "sandbox")]
mod sandbox {
    use std::os::raw::c_char;
    use std::ffi::CString;
    extern {
        fn sandbox(mode: *const c_char);
    }

    pub fn enter_sandbox(mode: &str) {
        unsafe {
            sandbox(CString::new(mode).unwrap().as_ptr());
        }
    }
}

#[cfg(not(feature = "sandbox"))]
mod sandbox {
    pub fn enter_sandbox(_mode: &str) {
        panic!("sandbox not available");
    }
}

pub fn main() {
    let args: Vec<String> = std::env::args().collect();
    let mut opts = Options::new();
    opts.optopt("", "sandbox", "sandbox mode", "MODE");
    let matches = opts.parse(&args[1..]).unwrap();

    let mut connection = unsafe {
        UnixStream::from_raw_fd(3)
    };

    let seed = {
        let mut random = unsafe { File::from_raw_fd(4) };
        random.read_i32::<NativeEndian>().unwrap()
    };

    let runner = Runner::new(seed).unwrap();
    let mut server = rmp_rpc::Server::new(runner);

    if let Some(mode) = matches.opt_str("sandbox") {
        sandbox::enter_sandbox(&mode);
    }

    server.serve(&mut connection);
}
