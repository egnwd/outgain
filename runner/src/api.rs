use mrusty::*;

use protocol::{WorldState, Entity};

mrusty_class!(WorldState, "WorldState", {
    def!("time", |mruby, slf: (&WorldState)| {
        mruby.fixnum(slf.time as i32)
    });

    def!("entities", |mruby, slf: (&WorldState)| {
        let entities = slf.entities.iter()
            .map(|e| mruby.obj(e.clone()))
            .collect();

        mruby.array(entities)
    });
});

mrusty_class!(Entity, "Entity", {
    def!("id", |mruby, slf: (&Entity)| {
        mruby.fixnum(slf.id as i32)
    });

    def!("x", |mruby, slf: (&Entity)| {
        mruby.float(slf.x)
    });

    def!("y", |mruby, slf: (&Entity)| {
        mruby.float(slf.y)
    });

    def!("radius", |mruby, slf: (&Entity)| {
        mruby.float(slf.radius)
    });

    def!("distance", |mruby, slf: (&Entity), other: (&Entity)| {
        let dx = slf.x - other.x;
        let dy = slf.y - other.y;

        return mruby.float(dx * dx + dy * dy);
    });
});
