use std::any::{Any, TypeId};
use std::sync::Arc;

use crate::object::Object;

pub struct Float {
    pub obj: Object,
}

impl Float {
    pub fn new(value: f64) -> Self {
        let mut obj = Object::new(Some(Box::new(value)),"float".to_string());

        let add_func: Arc<dyn Fn(&Box<dyn Any>) -> Object + Send + Sync> =
            Arc::new(move |other: &Box<dyn Any>| -> Object {
                let a = value;

                if let Some(b) = other.downcast_ref::<f64>() {
                    Float::new(a + *b).obj
                } else if let Some(b) = other.downcast_ref::<i64>() {
                    Float::new(a + *b as f64).obj
                } else {
                    panic!("Error!")
                }
            });

        obj.set("add".to_string(), add_func);

        Self { obj }
    }

    pub fn get_self_value(&self) -> Option<&f64> {
        self.obj.get_self_value()
    }
}
