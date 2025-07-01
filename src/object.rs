use std::fmt;
use std::any::Any;
use std::collections::HashMap;

pub struct Object {
    self_value: Option<Box<dyn Any>>,
    properties: HashMap<String, Box<dyn Any>>,
    class: String,
}

impl Object {
    pub fn new(self_value: Option<Box<dyn Any>>, class: String) -> Self {
        Self {
            self_value,
            properties: HashMap::new(),
            class,
        }
    }

    pub fn set_self_value(&mut self, value: Box<dyn Any>, class: String) {
        self.self_value = Some(value);
        self.class = class;
    }

    pub fn get_class(&self) -> &str {
        &self.class
    }

    pub fn get_self_value<T: 'static>(&self) -> Option<&T> {
        self.self_value.as_ref()?.downcast_ref::<T>()
    }

    pub fn set<T: 'static>(&mut self, key: String, value: T) {
        self.properties.insert(key, Box::new(value));
    }

    pub fn get<T: 'static>(&self, key: &str) -> Option<&T> {
        self.properties.get(key)?.downcast_ref::<T>()
    }

    pub fn properties_len(&self) -> usize {
        self.properties.len()
    }

    pub fn property_names(&self) -> Vec<String> {
        self.properties.keys().cloned().collect()
    }
}

impl fmt::Display for Object {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        match self.class.as_str() {
            "integer" => {
                if let Some(integer) = self.get_self_value::<i64>() {
                    write!(f, "{}", integer)
                } else {
                    write!(f, "Error!")
                }
            }
            "float" => {
                if let Some(float) = self.get_self_value::<f64>() {
                    write!(f, "{:.10}", float)
                } else {
                    write!(f, "Error!")
                }
            }
            _ => {
                if let Some(ptr) = self.self_value.as_ref() {
                    let addr = &**ptr as *const dyn Any as *const ();
                    write!(f, "<{} at {:p}>", self.class, addr)
                } else {
                    write!(f, "nil")
                }
            }
        }
    }
}
