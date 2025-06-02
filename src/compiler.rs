use crate::parser::Expr;
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
    }
}

fn compile_expr(expr: Expr) -> Vec<u8> {
    match expr {
        Expr::Lit(n) => match n {
            Value::IntVal(d) => {
                let mut bb = vec![OpCode::PushInt.as_byte()];
                let mut d_bytes = d.to_be_bytes().to_vec();
                bb.append(&mut d_bytes);
                bb
            }
        },
        Expr::BinOp(_, _, _) => todo!(),
    }
}

enum OpCode {
    _NOP,
    PushInt,
}

impl OpCode {
    fn as_byte(&self) -> u8 {
        match self {
            OpCode::_NOP => 0x00,
            OpCode::PushInt => 0x01,
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
