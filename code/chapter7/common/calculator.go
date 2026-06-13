package common

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"math"
	"strings"
)

// Calculator 是一个简单的数学计算工具
// 支持基本运算 (+, -, *, /)、取模 (%)、幂运算 (^) 和 sqrt 函数
type Calculator struct {
	Name        string
	Description string
}

// NewCalculator 创建计算器工具实例
func NewCalculator() *Calculator {
	return &Calculator{
		Name:        "my_calculator",
		Description: "简单的数学计算工具，支持基本运算(+,-,*,/)和sqrt函数",
	}
}

// Run 执行计算（满足 Tool 接口）
func (c *Calculator) Run(input string) string {
	return MyCalculate(input)
}

// MyCalculate 简单的数学计算函数
func MyCalculate(expression string) string {
	if strings.TrimSpace(expression) == "" {
		return "计算表达式不能为空"
	}

	fset := token.NewFileSet()
	expr, err := parser.ParseExprFrom(fset, "expr", []byte(expression), 0)
	if err != nil {
		return "计算失败，请检查表达式格式"
	}

	result, err := evalNode(expr)
	if err != nil {
		return "计算失败，请检查表达式格式"
	}

	return fmt.Sprintf("%v", result)
}

// evalNode 简化的表达式求值
func evalNode(node ast.Expr) (float64, error) {
	switch n := node.(type) {
	case *ast.BasicLit:
		var v float64
		_, err := fmt.Sscanf(n.Value, "%f", &v)
		if err != nil {
			return 0, err
		}
		return v, nil

	case *ast.BinaryExpr:
		left, err := evalNode(n.X)
		if err != nil {
			return 0, err
		}
		right, err := evalNode(n.Y)
		if err != nil {
			return 0, err
		}
		switch n.Op {
		case token.ADD:
			return left + right, nil
		case token.SUB:
			return left - right, nil
		case token.MUL:
			return left * right, nil
		case token.QUO:
			return left / right, nil
		case token.REM:
			return math.Mod(left, right), nil
		case token.XOR:
			return math.Pow(left, right), nil
		default:
			return 0, fmt.Errorf("不支持的运算符")
		}

	case *ast.CallExpr:
		funcName, ok := n.Fun.(*ast.Ident)
		if !ok {
			return 0, fmt.Errorf("不支持的函数调用")
		}
		if funcName.Name == "sqrt" {
			if len(n.Args) != 1 {
				return 0, fmt.Errorf("sqrt 函数需要一个参数")
			}
			arg, err := evalNode(n.Args[0])
			if err != nil {
				return 0, err
			}
			return math.Sqrt(arg), nil
		}
		return 0, fmt.Errorf("不支持的函数: %s", funcName.Name)

	case *ast.Ident:
		if n.Name == "pi" || n.Name == "Pi" {
			return math.Pi, nil
		}
		return 0, fmt.Errorf("未知标识符: %s", n.Name)

	case *ast.ParenExpr:
		return evalNode(n.X)

	case *ast.UnaryExpr:
		val, err := evalNode(n.X)
		if err != nil {
			return 0, err
		}
		switch n.Op {
		case token.SUB:
			return -val, nil
		case token.ADD:
			return val, nil
		default:
			return 0, fmt.Errorf("不支持的一元运算符")
		}

	default:
		return 0, fmt.Errorf("不支持的表达式类型")
	}
}
