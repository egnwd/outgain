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

impl Runner {
    pub fn new(source: String, seed: i32) -> Runner {
        let mrb = Mruby::new();

        let seed = mrb.fixnum(seed);
        let random = mrb.get_class("Random").unwrap().to_value();
        random.call("srand", vec![seed]).unwrap();

        mrb.filename("input");
        mrb.run(&source).unwrap();

        WorldState::require(mrb.clone());
        Entity::require(mrb.clone());

        Runner {
            mrb: mrb,
        }
    }

    pub fn tick(&mut self, request: TickRequest) -> TickResult {
        fn num_to_float(v: mrusty::Value) -> Result<f64, MrubyError> {
            v.to_f64().or_else(|_err| {
                v.to_i32().map(|x| x as f64)
            })
        }

        let object = self.mrb.top_self();
        let player = self.mrb.obj(request.player);
        let world = self.mrb.obj(request.world_state);
        let result = object.call("run", vec![player, world]).unwrap();

        result.to_vec().and_then(|v| {
            let dx = try!(num_to_float(v[0].clone()));
            let dy = try!(num_to_float(v[1].clone()));

            Ok(TickResult {
                dx: dx,
                dy: dy,
            })
        }).unwrap()
    }
}

impl rmp_rpc::Handler for Runner {
    fn request(&mut self, method: &str, params: Vec<rmp::Value>)
            -> Result<rmp::Value, rmp::Value> {

        let param = params.into_iter().next();
        match (method, param) {
            ("Runner.Tick", Some(param)) => {
                let request = rmp_serde::from_value(param).unwrap();

                let result = self.tick(request);

                let response = rmp_serde::to_value(&result);
                Ok(response)
            }

            _ => {
                Err(rmp::Value::String(format!("invalid RPC method call: {:?}", method)))
            }
        }
    }
}
