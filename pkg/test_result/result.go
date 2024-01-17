package test_result

type Result struct {
	Total           int
	Subtotal_1      int
	Subtotal_2      int
	Subtotal_3      int
	Subtotal_4      int
}

func MakeResult(answers [25]bool) Result {
	var sub1, sub2, sub3, sub4, tot int

	for i, v := range answers {
		if v == true {
			switch i {
			case 5, 10, 16, 19, 21, 22:
				sub1++
			case 3, 4, 12, 14, 20, 23:
				sub2++
			case 1, 7, 9, 11, 18, 24:
				sub3++
			case 2, 6, 8, 13, 15, 17:
				sub4++
			}
			tot++
		}
	}

	res := Result{
		Total:           tot,
		Subtotal_1:      sub1,
		Subtotal_2:      sub2,
		Subtotal_3:      sub3,
		Subtotal_4:      sub4,
	}

	return res
}

func calculateResult(res int) float32 {
	var subTotal float32
	subTotal = (float32(res) / 6) * 100
	return subTotal
}

func calculateTotalResult(res int) float32 {
	var TotalResult float32
	TotalResult = float32((float32(res) / 24) * 100)
	return TotalResult
}

func (r *Result) GetChooseSubTotal() []string {
	res := [4]float32{r.GetSubtotal_1(), r.GetSubtotal_2(), r.GetSubtotal_3(), r.GetSubtotal_4()}
	var per float32
	for _, number := range res {
		if number > per {
			per = number
		}
	}
	m := make([]string,0,4)

	if per == 0{
		m = append(m, "Try again")
		return m
	}
	if per == r.GetSubtotal_1() {
		m = append(m, "Your result is Subtotal_1")
	}
	if per == r.GetSubtotal_2() {
		m = append(m, "Your result is Subtotal_2")
	}
	if per == r.GetSubtotal_3() {
		m = append(m, "Your result is Subtotal_3")
	}
	if per == r.GetSubtotal_4() {
		m = append(m, "Your result is Subtotal_4")
	}
	return m
}

func (r *Result) GetSubtotal_1() float32 {
	return calculateResult(r.Subtotal_1)
}

func (r *Result) GetSubtotal_2() float32 {
	return calculateResult(r.Subtotal_2)
}

func (r *Result) GetSubtotal_3() float32 {
	return calculateResult(r.Subtotal_3)
}

func (r *Result) GetSubtotal_4() float32 {
	return calculateResult(r.Subtotal_4)
}

func (r *Result) GetTotal() float32 {
	return calculateTotalResult(r.Total)
}
