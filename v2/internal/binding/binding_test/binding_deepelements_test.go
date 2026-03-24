package binding_test

// Issues 2303, 3442, 3709

type DeepMessage struct {
	Msg string
}

type DeepElements struct {
	Single     []int
	Double     [][]string
	FourDouble [4][]float64
	DoubleFour [][4]int64
	Triple     [][][]int

	SingleMap      map[string]int
	SliceMap       map[string][]int
	DoubleSliceMap map[string][][]int

	ArrayMap        map[string][4]int
	DoubleArrayMap1 map[string][4][]int
	DoubleArrayMap2 map[string][][4]int
	DoubleArrayMap3 map[string][4][4]int

	OneStructs      []*DeepMessage
	TwoStructs      [3][]*DeepMessage
	ThreeStructs    [][][]DeepMessage
	MapStructs      map[string][]*DeepMessage
	MapTwoStructs   map[string][4][]DeepMessage
	MapThreeStructs map[string][][7][]*DeepMessage
}

func (x DeepElements) Get() DeepElements {
	return x
}

var DeepElementsTest = BindingTest{
	name: "DeepElements",
	structs: []interface{}{
		&DeepElements{},
	},
	exemptions:  nil,
	shouldError: false,
	want: `
export namespace binding_test {
            
                export class DeepMessage {
                    Msg: string;
            
                    static createFrom(source: any = {}) {
                        return new DeepMessage(source);
                    }
            
                    constructor(source: any = {}) {
                        if ('string' === typeof source) source = JSON.parse(source);
                        this.Msg = source["Msg"];
                    }
                }
                export class DeepElements {
                    Single: number[];
                    Double: string[][];
                    FourDouble: number[][];
                    DoubleFour: number[][];
                    Triple: number[][][];
                    SingleMap: Record<string, number>;
                    SliceMap: Record<string, Array<number>>;
                    DoubleSliceMap: Record<string, Array<Array<number>>>;
                    ArrayMap: Record<string, Array<number>>;
                    DoubleArrayMap1: Record<string, Array<Array<number>>>;
                    DoubleArrayMap2: Record<string, Array<Array<number>>>;
                    DoubleArrayMap3: Record<string, Array<Array<number>>>;
                    OneStructs: DeepMessage[];
                    TwoStructs: DeepMessage[][];
                    ThreeStructs: DeepMessage[][][];
                    MapStructs: Record<string, Array<DeepMessage>>;
                    MapTwoStructs: Record<string, Array<Array<DeepMessage>>>;
                    MapThreeStructs: Record<string, Array<Array<Array<DeepMessage>>>>;
            
                    static createFrom(source: any = {}) {
                        return new DeepElements(source);
                    }
            
                    constructor(source: any = {}) {
                        if ('string' === typeof source) source = JSON.parse(source);
                        this.Single = source["Single"];
                        this.Double = source["Double"];
                        this.FourDouble = source["FourDouble"];
                        this.DoubleFour = source["DoubleFour"];
                        this.Triple = source["Triple"];
                        this.SingleMap = source["SingleMap"];
                        this.SliceMap = source["SliceMap"];
                        this.DoubleSliceMap = source["DoubleSliceMap"];
                        this.ArrayMap = source["ArrayMap"];
                        this.DoubleArrayMap1 = source["DoubleArrayMap1"];
                        this.DoubleArrayMap2 = source["DoubleArrayMap2"];
                        this.DoubleArrayMap3 = source["DoubleArrayMap3"];
                        this.OneStructs = this.convertValues(source["OneStructs"], DeepMessage);
                        this.TwoStructs = this.convertValues(source["TwoStructs"], DeepMessage);
                        this.ThreeStructs = this.convertValues(source["ThreeStructs"], DeepMessage);
                        this.MapStructs = this.convertValues(source["MapStructs"], Array<DeepMessage>, true);
                        this.MapTwoStructs = this.convertValues(source["MapTwoStructs"], Array<Array<DeepMessage>>, true);
                        this.MapThreeStructs = this.convertValues(source["MapThreeStructs"], Array<Array<Array<DeepMessage>>>, true);
                    }
            
                        convertValues(a: any, classs: any, asMap: boolean = false): any {
                            if (!a) {
                                return a;
                            }
                            if (a.slice && a.map) {
                                return (a as any[]).map(elem => this.convertValues(elem, classs));
                            } else if ("object" === typeof a) {
                                if (asMap) {
                                    for (const key of Object.keys(a)) {
                                        a[key] = new classs(a[key]);
                                    }
                                    return a;
                                }
                                return new classs(a);
                            }
                            return a;
                        }
                }
            
            }
`,
}
