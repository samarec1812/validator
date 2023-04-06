package validator

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		checkErr func(err error) bool
	}{
		{
			name: "invalid struct: interface",
			args: args{
				v: new(any),
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrNotStruct)
			},
		},
		{
			name: "invalid struct: map",
			args: args{
				v: map[string]string{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrNotStruct)
			},
		},
		{
			name: "invalid struct: string",
			args: args{
				v: "some string",
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrNotStruct)
			},
		},
		{
			name: "valid struct with no fields",
			args: args{
				v: struct{}{},
			},
			wantErr: false,
		},
		{
			name: "valid struct with untagged fields",
			args: args{
				v: struct {
					f1 string
					f2 string
				}{},
			},
			wantErr: false,
		},
		{
			name: "valid struct with unexported fields",
			args: args{
				v: struct {
					foo string `validate:"len:10"`
				}{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				e := &ValidationErrors{}
				return errors.As(err, e) && e.Error() == ErrValidateForUnexportedFields.Error()
			},
		},
		{
			name: "invalid validator syntax",
			args: args{
				v: struct {
					Foo string `validate:"len:abcdef"`
				}{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				e := &ValidationErrors{}
				return errors.As(err, e) && e.Error() == ErrInvalidValidatorSyntax.Error()
			},
		},
		{
			name: "valid struct with tagged fields",
			args: args{
				v: struct {
					Len       string `validate:"len:20"`
					LenZ      string `validate:"len:0"`
					InInt     int    `validate:"in:20,25,30"`
					InNeg     int    `validate:"in:-20,-25,-30"`
					InStr     string `validate:"in:foo,bar"`
					MinInt    int    `validate:"min:10"`
					MinIntNeg int    `validate:"min:-10"`
					MinStr    string `validate:"min:10"`
					MinStrNeg string `validate:"min:-1"`
					MaxInt    int    `validate:"max:20"`
					MaxIntNeg int    `validate:"max:-2"`
					MaxStr    string `validate:"max:20"`
				}{
					Len:       "abcdefghjklmopqrstvu",
					LenZ:      "",
					InInt:     25,
					InNeg:     -25,
					InStr:     "bar",
					MinInt:    15,
					MinIntNeg: -9,
					MinStr:    "abcdefghjkl",
					MinStrNeg: "abc",
					MaxInt:    16,
					MaxIntNeg: -3,
					MaxStr:    "abcdefghjklmopqrst",
				},
			},
			wantErr: false,
		},
		{
			name: "wrong length",
			args: args{
				v: struct {
					Lower    string `validate:"len:24"`
					Higher   string `validate:"len:5"`
					Zero     string `validate:"len:3"`
					BadSpec  string `validate:"len:%12"`
					Negative string `validate:"len:-6"`
				}{
					Lower:    "abcdef",
					Higher:   "abcdef",
					Zero:     "",
					BadSpec:  "abc",
					Negative: "abcd",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong in",
			args: args{
				v: struct {
					InA     string `validate:"in:ab,cd"`
					InB     string `validate:"in:aa,bb,cd,ee"`
					InC     int    `validate:"in:-1,-3,5,7"`
					InD     int    `validate:"in:5-"`
					InEmpty string `validate:"in:"`
				}{
					InA:     "ef",
					InB:     "ab",
					InC:     2,
					InD:     12,
					InEmpty: "",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong min",
			args: args{
				v: struct {
					MinA string `validate:"min:12"`
					MinB int    `validate:"min:-12"`
					MinC int    `validate:"min:5-"`
					MinD int    `validate:"min:"`
					MinE string `validate:"min:"`
				}{
					MinA: "ef",
					MinB: -22,
					MinC: 12,
					MinD: 11,
					MinE: "abc",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong max",
			args: args{
				v: struct {
					MaxA string `validate:"max:2"`
					MaxB string `validate:"max:-7"`
					MaxC int    `validate:"max:-12"`
					MaxD int    `validate:"max:5-"`
					MaxE int    `validate:"max:"`
					MaxF string `validate:"max:"`
				}{
					MaxA: "efgh",
					MaxB: "ab",
					MaxC: 22,
					MaxD: 12,
					MaxE: 11,
					MaxF: "abc",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 6)
				return true
			},
		},
		{
			name: "valid struct with untagged fields slice",
			args: args{
				v: struct {
					f1 string
					f2 []string
				}{},
			},
			wantErr: false,
		},
		{
			name: "valid struct with unexported fields slice",
			args: args{
				v: struct {
					foo []string `validate:"len:10"`
					bar []int
				}{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				e := &ValidationErrors{}
				return errors.As(err, e) && e.Error() == ErrValidateForUnexportedFields.Error()
			},
		},
		{
			name: "valid struct with tagged fields with slice",
			args: args{
				v: struct {
					Len            string   `validate:"len:20"`
					LenSlice       []string `validate:"len:5"`
					LenZ           string   `validate:"len:0"`
					LenSliceZ      []string `validate:"len:0"`
					InInt          int      `validate:"in:20,25,30"`
					InSliceInt     []int    `validate:"in:4,10,14"`
					InNeg          int      `validate:"in:-20,-25,-30"`
					InSliceNeg     []int    `validate:"in:-3,-13,-10"`
					InStr          string   `validate:"in:foo,bar"`
					InSliceStr     []string `validate:"in:foo,bar,pon,legal"`
					MinInt         int      `validate:"min:10"`
					MinSliceInt    []int    `validate:"min:13"`
					MinIntNeg      int      `validate:"min:-10"`
					MinSliceNeg    []int    `validate:"min:-5"`
					MinStr         string   `validate:"min:10"`
					MinSliceStr    []string `validate:"min:9"`
					MinStrNeg      string   `validate:"min:-1"`
					MinSliceStrNeg []string `validate:"min:-4"`
					MaxInt         int      `validate:"max:20"`
					MaxSliceInt    []int    `validate:"max:15"`
					MaxIntNeg      int      `validate:"max:-2"`
					MaxSliceNeg    []int    `validate:"max:-6"`
					MaxStr         string   `validate:"max:20"`
					MaxSliceStr    []string `validate:"max:15"`
					MaxSliceStrZ   []string `validate:"max:0"`
					MaxSliceStrNeg []string `validate:"max:-4"`
				}{
					Len:            "abcdefghjklmopqrqwer",
					LenSlice:       []string{"abcde", "bfert", "kolpo", "piton"},
					LenZ:           "",
					LenSliceZ:      []string{"", "", ""},
					InInt:          25,
					InSliceInt:     []int{4, 4, 14, 10},
					InNeg:          -25,
					InSliceNeg:     []int{-10, -13, -13, -10, -3},
					InStr:          "bar",
					InSliceStr:     []string{"legal", "pon", "foo"},
					MinInt:         15,
					MinSliceInt:    []int{14, 16, 13, 99},
					MinIntNeg:      -9,
					MinSliceNeg:    []int{-1, -5, 0, 15},
					MinStr:         "abcdefghjkl",
					MinSliceStr:    []string{"akdwaobdnoai", "japdojao[djaoadmwpa[dm", "apdwoandoadomno", "adwkolpok"},
					MinStrNeg:      "abc",
					MinSliceStrNeg: []string{"awd", "", "ajpdjapwdapowdjapo"},
					MaxInt:         16,
					MaxSliceInt:    []int{1, 0, -5, 15, 10, 12, -34},
					MaxIntNeg:      -3,
					MaxSliceNeg:    []int{-10, -6, -40, -14},
					MaxStr:         "abcdefghjklmopqrst",
					MaxSliceStr:    []string{"awd", "", "apdmpaodm", "awdftghujkolawd"},
					MaxSliceStrZ:   []string{"", "", "", "", ""},
					MaxSliceStrNeg: []string{},
				},
			},
			wantErr: false,
		},
		{
			name: "wrong slice length",
			args: args{
				v: struct {
					Lower    []string `validate:"len:24"`
					Higher   []string `validate:"len:5"`
					Zero     []string `validate:"len:3"`
					BadSpec  []string `validate:"len:%12"`
					Negative []string `validate:"len:-6"`
				}{
					Lower:    []string{"awddwaawdawddwaawdsa", "abcdef"},
					Higher:   []string{"abcdef", "awdajdpajpd", "apindwpajdpwaodjapd"},
					Zero:     []string{"", "", ""},
					BadSpec:  []string{"abc", "adwadad"},
					Negative: []string{"abcd", "a", "", "apkwd[pad"},
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong slice in",
			args: args{
				v: struct {
					InA     []string `validate:"in:ab,cd"`
					InB     []string `validate:"in:aa,bb,cd,ee"`
					InC     []int    `validate:"in:-1,-3,5,7"`
					InD     []int    `validate:"in:5-"`
					InEmpty []string `validate:"in:"`
				}{
					InA:     []string{"ab", "ef", "cd"},
					InB:     []string{"ab", "bb"},
					InC:     []int{-3, 2},
					InD:     []int{12, 0, 9, 4},
					InEmpty: []string{""},
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong slice min",
			args: args{
				v: struct {
					MinA []string `validate:"min:12"`
					MinB []int    `validate:"min:-12"`
					MinC []int    `validate:"min:5"`
					MinD []int    `validate:"min:-3"`
					MinE []string `validate:"min:"`
				}{
					MinA: []string{"ef", "", "adwadaaojwkdpoadd"},
					MinB: []int{-22, 0, 24},
					MinC: []int{12, 3, 0},
					MinD: []int{11, 0, 3, -90},
					MinE: []string{},
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong slice max",
			args: args{
				v: struct {
					MaxA []string `validate:"max:2"`
					MaxB []string `validate:"max:-7"`
					MaxC []int    `validate:"max:-12"`
					MaxD []int    `validate:"max:5"`
					MaxE []int    `validate:"max:3"`
					MaxF []string `validate:"max:2"`
				}{
					MaxA: []string{"efgh", "adw", "a"},
					MaxB: []string{"ab", "adwad"},
					MaxC: []int{22, 0, 12},
					MaxD: []int{12, 2, 5},
					MaxE: []int{-9, 10, 0},
					MaxF: []string{"ab", "a", "b"},
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.args.v)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, tt.checkErr(err), "test expect an error, but got wrong error type")
			} else {
				assert.NoError(t, err)
			}
		})
	}

}
