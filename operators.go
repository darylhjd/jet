package jet

import "errors"

//----------- Logical operators ---------------//

// NOT returns negation of bool expression result
func NOT(exp BoolExpression) BoolExpression {
	return newPrefixBoolOperator(exp, "NOT")
}

// BIT_NOT inverts every bit in integer expression result
func BIT_NOT(expr IntegerExpression) IntegerExpression {
	return newPrefixIntegerOperator(expr, "~")
}

//----------- Comparison operators ---------------//

// EXISTS checks for existence of the rows in subQuery
func EXISTS(subQuery SelectStatement) BoolExpression {
	return newPrefixBoolOperator(subQuery, "EXISTS")
}

// Returns a representation of "a=b"
func eq(lhs, rhs Expression) BoolExpression {
	return newBinaryBoolOperator(lhs, rhs, "=")
}

// Returns a representation of "a!=b"
func notEq(lhs, rhs Expression) BoolExpression {
	return newBinaryBoolOperator(lhs, rhs, "!=")
}

func isDistinctFrom(lhs, rhs Expression) BoolExpression {
	return newBinaryBoolOperator(lhs, rhs, "IS DISTINCT FROM")
}

func isNotDistinctFrom(lhs, rhs Expression) BoolExpression {
	return newBinaryBoolOperator(lhs, rhs, "IS NOT DISTINCT FROM")
}

// Returns a representation of "a<b"
func lt(lhs Expression, rhs Expression) BoolExpression {
	return newBinaryBoolOperator(lhs, rhs, "<")
}

// Returns a representation of "a<=b"
func ltEq(lhs, rhs Expression) BoolExpression {
	return newBinaryBoolOperator(lhs, rhs, "<=")
}

// Returns a representation of "a>b"
func gt(lhs, rhs Expression) BoolExpression {
	return newBinaryBoolOperator(lhs, rhs, ">")
}

// Returns a representation of "a>=b"
func gtEq(lhs, rhs Expression) BoolExpression {
	return newBinaryBoolOperator(lhs, rhs, ">=")
}

// --------------- CASE operator -------------------//

// CaseOperator is interface for SQL case operator
type CaseOperator interface {
	Expression

	WHEN(condition Expression) CaseOperator
	THEN(then Expression) CaseOperator
	ELSE(els Expression) CaseOperator
}

type caseOperatorImpl struct {
	expressionInterfaceImpl

	expression Expression
	when       []Expression
	then       []Expression
	els        Expression
}

// CASE create CASE operator with optional list of expressions
func CASE(expression ...Expression) CaseOperator {
	caseExp := &caseOperatorImpl{}

	if len(expression) > 0 {
		caseExp.expression = expression[0]
	}

	caseExp.expressionInterfaceImpl.parent = caseExp

	return caseExp
}

func (c *caseOperatorImpl) WHEN(when Expression) CaseOperator {
	c.when = append(c.when, when)
	return c
}

func (c *caseOperatorImpl) THEN(then Expression) CaseOperator {
	c.then = append(c.then, then)
	return c
}

func (c *caseOperatorImpl) ELSE(els Expression) CaseOperator {
	c.els = els

	return c
}

func (c *caseOperatorImpl) serialize(statement statementType, out *sqlBuilder, options ...serializeOption) error {
	if c == nil {
		return errors.New("jet: Case Expression is nil. ")
	}

	out.writeString("(CASE")

	if c.expression != nil {
		err := c.expression.serialize(statement, out)

		if err != nil {
			return err
		}
	}

	if len(c.when) == 0 || len(c.then) == 0 {
		return errors.New("jet: Invalid case Statement. There should be at least one when/then Expression pair. ")
	}

	if len(c.when) != len(c.then) {
		return errors.New("jet: When and then Expression count mismatch. ")
	}

	for i, when := range c.when {
		out.writeString("WHEN")
		err := when.serialize(statement, out, noWrap)

		if err != nil {
			return err
		}

		out.writeString("THEN")
		err = c.then[i].serialize(statement, out, noWrap)

		if err != nil {
			return err
		}
	}

	if c.els != nil {
		out.writeString("ELSE")
		err := c.els.serialize(statement, out, noWrap)

		if err != nil {
			return err
		}
	}

	out.writeString("END)")

	return nil
}
