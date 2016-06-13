extern crate gcc;
extern crate pkg_config;

fn main() {
    if std::env::var_os("CARGO_FEATURE_SANDBOX").is_some() {
        pkg_config::find_library("libseccomp").unwrap();
        gcc::compile_library("libsandbox.a", &["src/sandbox.c"]);
    }
}
