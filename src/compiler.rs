
use crate::parser::Expr;
use crate::parser::Operator;
use crate::parser::Stmt;
use crate::parser::Value;

pub fn compile(stmts: Vec<Stmt>) -> Vec<u8> {
    let mut bbytes = Vec::new();

    for stmt in stmts.into_iter() {
        bbytes.append(&mut compile_stmt(stmt));
    }

    bbytes
}

fn compile_stmt(stmt: Stmt) -> Vec<u8> {
    match stmt {
        Stmt::ExprStmt(expr) => compile_expr(expr),
        Stmt::Print(expr) => {
            let mut v = compile_expr(expr);
            v.push(OpCode::PrintTop.into());
            v
        }
    }
}

fn compile_expr(expr: Expr) -> Vec<u8> {
    match expr {
        Expr::Lit(n) => match n {
            Value::IntVal(d) => {
                let mut bb = vec![OpCode::PushInt.into()];
                let mut d_bytes = d.to_be_bytes().to_vec();
                bb.append(&mut d_bytes);
                bb
            }
        },
        Expr::BinOp(op, left, right) => {
            let mut bb = vec![];
            let mut leftbb = compile_expr(*left);
            let mut rightbb = compile_expr(*right);
            let opcode: OpCode = op.into();
            let opcodebb: u8 = opcode.into();
            bb.append(&mut leftbb);
            bb.append(&mut rightbb);
            bb.push(opcodebb);
            bb
        },
    }
}

pub enum OpCode {
    _NOP,
    PushInt,
    PrintTop,
    Add,
    Sub,
    Mul,
    Div,
}


impl From<Operator> for OpCode {
    fn from(value: Operator) -> Self {
        match value {
            Operator::ADD => OpCode::Add,
            Operator::SUB => OpCode::Sub,
            Operator::MUL => OpCode::Mul,
            Operator::DIV => OpCode::Div,
        }
    }
}

impl From<OpCode> for u8 {
    fn from(value: OpCode) -> Self {
        match value {
            OpCode::_NOP => 0x00,
            OpCode::PushInt => 0x01,
            OpCode::PrintTop => 0x02,
            OpCode::Add => 0x03,
            OpCode::Sub => 0x04,
            OpCode::Mul => 0x05,
            OpCode::Div => 0x06,
        }
    }
}

impl TryFrom<u8> for OpCode {
    type Error = String;

    fn try_from(value: u8) -> Result<Self, Self::Error> {
        match value {
          0 => Ok(OpCode::_NOP),
          1 => Ok(OpCode::PushInt),
          2 => Ok(OpCode::PrintTop),
          3 => Ok(OpCode::Add),
          4 => Ok(OpCode::Sub),
          5 => Ok(OpCode::Mul),
          6 => Ok(OpCode::Div),
          x => Err(format!("cannot convert '{}' to an OpCode", x)),
        }
    }
}

#[cfg(test)]
mod tests {
    // Note this useful idiom: importing names from outer (for mod tests) scope.
    use super::*;

    #[test]
    fn test_intval() {
        let input = vec![Stmt::ExprStmt(Expr::Lit(Value::IntVal(44)))];

        let want: Vec<u8> = vec![0x01, 0, 0, 0, 0, 0, 0, 0, 44];

        let got = compile(input);

        assert_eq!(want, got);
    }
}
