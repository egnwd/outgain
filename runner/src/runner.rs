use mrusty::{
    self,
    Mruby,
    MrubyType,
    MrubyImpl,
    MrubyError,
    MrubyFile
};

use rmp;
use rmp_rpc;
use rmp_serde;

use protocol::{
    TickRequest,
    TickResult,
    WorldState,
    Entity,
};

pub struct Runner {
    mrb: MrubyType,
}

fn num_to_float(v: mrusty::Value) -> Result<f64, MrubyError> {
    v.to_i32().map(|x| x as f64).or_else(|_| v.to_f64())
}

impl Runner {
    pub fn new(seed: i32) -> Result<Runner, MrubyError> {
        let mrb = Mruby::new();

        let seed = mrb.fixnum(seed);
        let random = try!(mrb.get_class("Random")).to_value();
        try!(random.call("srand", vec![seed]));

        WorldState::require(mrb.clone());
        Entity::require(mrb.clone());

        Ok(Runner {
            mrb: mrb,
        })
    }

    pub fn tick(&mut self, request: TickRequest) -> Result<TickResult, MrubyError> {
        let object = self.mrb.top_self();
        let player = self.mrb.obj(request.player);
        let world = self.mrb.obj(request.world_state);
        let result = try!(object.call("run", vec![player, world]));
        let result = try!(result.to_vec());

        let dx = try!(num_to_float(result[0].clone()));
        let dy = try!(num_to_float(result[1].clone()));

        Ok(TickResult {
            dx: dx,
            dy: dy,
        })
    }

    pub fn load(&mut self, source: String) -> Result<(), MrubyError> {
        self.mrb.filename("input");
        try!(self.mrb.run(&source));

        Ok(())
    }
}

macro_rules! rpc_try {
    ($e:expr) => {
        match $e {
            Ok(v) => v,
            Err(e) => {
                return Err(::rmp::Value::String(e.to_string()))
            }
        }
    }
}

impl rmp_rpc::Handler for Runner {
    fn request(&mut self, method: &str, params: Vec<rmp::Value>)
            -> Result<rmp::Value, rmp::Value> {

        let param = params.into_iter().next();
        match (method, param) {
            ("Runner.Load", Some(param)) => {
                let request = rpc_try!(rmp_serde::from_value(param));

                rpc_try!(self.load(request));

                Ok(rmp::Value::Nil)
            }

            ("Runner.Tick", Some(param)) => {
                let request = rpc_try!(rmp_serde::from_value(param));

                let result = rpc_try!(self.tick(request));

                Ok(rmp_serde::to_value(&result))
            }

            _ => {
                Err(rmp::Value::String(format!("invalid RPC method call: {:?}", method)))
            }
        }
    }
}
