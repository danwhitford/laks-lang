use std::{io::Error, io::Write};

use crate::{compiler::OpCode, parser::Value};

struct Vm<W: Write> {
    chunk: Vec<u8>,
    ip: usize,
    out: W,
    val_stack: Vec<Value>,
}

pub fn run(chunk: Vec<u8>, out: &mut impl Write) {
    let mut vm = Vm {
        chunk,
        ip: 0,
        out,
        val_stack: Vec::new(),
    };
    vm.run_chunk()
}

impl<W: Write> Vm<W> {
    fn run_chunk(&mut self) {
        while let Some(code) = self.read() {
            match OpCode::try_from(*code) {
                Ok(opcode) => match opcode {
                    OpCode::_NOP => (),
                    OpCode::PushInt => {
                        let val = self.read_int();
                        self.val_stack.push(val);
                    }
                    OpCode::PrintTop => {
                        let res = self.print_val();
                        if res.is_err() {
                            panic!("failed to write. {:?}", res.err());
                        }
                    }
                    OpCode::Add => {
                        let a = self.val_stack.pop().expect("stack is empty");
                        let b = self.val_stack.pop().expect("stack is empty");
                        match (a, b) {
                            (Value::IntVal(x), Value::IntVal(y)) => self.val_stack.push(Value::IntVal(x + y)),
                        }
                    },
                    OpCode::Sub => todo!(),
                    OpCode::Mul => todo!(),
                    OpCode::Div => todo!(),
                },
                Err(err) => panic!("op code not recognised. {}", err),
            }
        }
    }

    fn read(&mut self) -> Option<&u8> {
        let b = self.chunk.get(self.ip);
        self.ip += 1;
        b
    }

    fn read_int(&mut self) -> Value {
        let intbytes = &self.chunk[self.ip..self.ip + 8];
        self.ip += 8;
        Value::IntVal(i64::from_be_bytes(intbytes.try_into().unwrap()))
    }

    fn print_val(&mut self) -> Result<(), Error> {
        match self.val_stack.pop().expect("stack empty") {
            Value::IntVal(d) => write!(self.out, "{}\n", d),
        }
    }
}

#[cfg(test)]
mod tests {
    // Note this useful idiom: importing names from outer (for mod tests) scope.
    use super::*;

    #[test]
    fn test_intval() {
        let input: Vec<u8> = vec![0x01, 0, 0, 0, 0, 0, 0, 0, 44, 0x02];

        let want = "44\n";

        let mut buf = vec![];
        run(input, &mut buf);
        let got = String::from_utf8(buf).expect("should be valid utf8 string");

        assert_eq!(want, got);
    }
}
