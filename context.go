package hpack

import (
	"flag"
	"fmt"
	. "github.com/jxck/color"
	. "github.com/jxck/logger"
	"log"
)

var verbose bool
var loglevel int

// Request or Response
type CXT bool

const (
	REQUEST  CXT = true
	RESPONSE     = false
)

func init() {
	flag.BoolVar(&verbose, "v", false, "verbose out")
	flag.IntVar(&loglevel, "l", 0, "log level (1 ERR, 2 WARNING, 3 INFO, 4 DEBUG)")
	flag.Parse()
	LogLevel(loglevel)
	Verbose(verbose)
	log.SetFlags(log.Lshortfile)
}

type Context struct {
	CXT
	HT *HeaderTable
	RS *ReferenceSet
	ES *EmittedSet
}

func NewContext(context CXT, SETTINGS_HEADER_TABLE_SIZE int) Context {
	return Context{
		HT:  NewHeaderTable(SETTINGS_HEADER_TABLE_SIZE),
		RS:  NewReferenceSet(),
		ES:  NewEmittedSet(),
		CXT: context,
	}
}

func (c *Context) Decode(wire []byte) {
	// 各デコードごとに前回のをリセットする。
	c.ES = NewEmittedSet()
	c.RS.Reset()
	Info(Red("clean Emitted Set"))
	Debug(Cyan(
		"\n===== Before Decode =====")+
		"%v"+Cyan(
		"==========================="),
		c.Dump())

	frames := Decode(wire, c.CXT)
	for _, frame := range frames {
		switch f := frame.(type) {
		case *IndexedHeader:
			index := int(f.Index)

			if index == 0 {
				/**
				 * idx=0 の場合 Reference Set を空にする
				 */
				Info(Red("Empty ReferenceSet"))
				c.RS.Empty()
				continue
			}

			var headerField *HeaderField

			if index > c.HT.Len() {
				/**
				 * Static Header Table の中にある場合
				 */
				i := index - c.HT.Len() - 1
				headerField = StaticHeaderTable[i]

				if c.RS.Has(headerField) {
					/**
					 * 参照が Reference Set にあった場合
					 * Reference Set から消す
					 */

					Info(Red(fmt.Sprintf("== Indexed - Remove ==")))
					Info(fmt.Sprintf("  idx = %v", index))
					Info(fmt.Sprintf("  -> HT[%v] = %v", index, headerField))

					// remove
					c.RS.Remove(headerField)
				} else {
					/**
					* 参照が Reference Set に無い場合
					* 該当のエントリを取り出す
					 */

					Info(Red("== Indexed - Add =="))
					Info(fmt.Sprintf("  idx = %v", index))
					Info(fmt.Sprintf("  -> ST[%v] = %v", index, headerField))

					// Emit
					Info(Blue("\tEmit"))
					c.ES.Emit(headerField)

					// ヘッダテーブルにコピーする
					Info(Blue("\tAdd to HT"))
					c.Push(headerField)

					// その参照を RefSet に追加する
					Info(Blue("\tAdd to RS"))
					c.RS.Add(headerField, EMITTED)
				}
			} else {
				/**
				 * Header Table の中にある場合
				 */

				// 実態は配列なので 0 オリジン
				i := index - 1
				headerField = c.HT.HeaderFields[i]

				if c.RS.Has(headerField) {
					/**
					 * 参照が Reference Set にあった場合
					 * Reference Set から消す
					 */
					Info(Red(fmt.Sprintf("== Indexed - Remove ==")))
					Info(fmt.Sprintf("  idx = %v", index))
					Info(fmt.Sprintf("  -> HT[%v] = %v", index, headerField))

					// remove
					c.RS.Remove(headerField)
				} else {
					/**
					* 参照が Reference Set に無い場合
					 */
					Info(Red("== Indexed - Add =="))
					Info(fmt.Sprintf("  idx = %v", index))
					Info(fmt.Sprintf("  -> HT[%v] = %v", index, headerField))

					// Emit
					Info(Blue("\tEmit"))
					c.ES.Emit(headerField)

					// その参照を RefSet に追加する
					Info(Blue("\tAdd to RS"))
					c.RS.Add(headerField, EMITTED)
				}
			}
		case *IndexedLiteral:

			// Index 先の Name と Literal Value から HeaderField を生成
			index := int(f.Index)
			var name, value string

			if index > c.HT.Len() {
				/**
				 * Static Header Table の中にある場合
				 */
				i := index - c.HT.Len() - 1
				name = StaticHeaderTable[i].Name
			} else {
				/**
				 * Header Table の中にある場合
				 */
				i := index - 1
				name = c.HT.HeaderFields[i].Name
			}

			value = f.ValueString

			// Header Field 生成
			headerField := NewHeaderField(name, value)

			Info(Red("== Indexed Literal =="))
			Info(fmt.Sprintf("  Indexed name (idx = %v)", index))
			Info(fmt.Sprintf("  -> ST[%v].Name = %v", index, name))
			Info(fmt.Sprintf("  Literal value (len = %v)", f.ValueLength))
			Info(fmt.Sprintf("  %v", f.ValueString))

			if f.Indexing {
				/**
				 * HT に追加する場合
				 */

				// Emit
				Info(Blue("\tEmit"))
				c.ES.Emit(headerField)

				// ヘッダテーブルにコピーする
				Info(Blue("\tAdd to HT"))
				c.Push(headerField)

				// その参照を RefSet に追加する
				Info(Blue("\tAdd to RS"))
				c.RS.Add(headerField, EMITTED)

			} else {
				/**
				 * HT に追加しない場合
				 */

				// Emit
				Info(Blue("\tEmit"))
				c.ES.Emit(headerField)
			}

		case *StringLiteral:
			Info(Red(fmt.Sprintf("== String Literal (%v) ==", f)))

			headerField := NewHeaderField(f.NameString, f.ValueString)
			if f.Indexing {
				// HT に追加する場合

				// Emit
				Info(Blue("\tEmit"))
				c.ES.Emit(headerField)

				// ヘッダテーブルにコピーする
				Info(Blue("\tAdd to HT"))
				c.Push(headerField)

				// その参照を RefSet に追加する
				Info(Blue("\tAdd to RS"))
				c.RS.Add(headerField, EMITTED)

			} else {
				// HT に追加しない場合

				// Emit
				Info(Blue("\tEmit"))
				c.ES.Emit(headerField)
			}

		default:
			log.Fatal("%T", f)
		}
	}

	// reference set の残りを全て emit する
	for _, referencedField := range *c.RS {
		if !referencedField.Emitted {
			headerField := referencedField.HeaderField
			Info(Blue("\tEmit rest entries ")+"%v", *headerField)
			c.ES.Emit(headerField)
		}
	}
}

// removing entry from top
// until make space of size in Header Table
// Evict された参照を RS からも消すために、 Context の方でやる。
func (c *Context) Eviction() {
	for c.HT.Size() > c.HT.HEADER_TABLE_SIZE {
		// サイズが収まるまで減らす
		Info(Red("Eviction")+" %v", c.HT.HeaderFields[len(c.HT.HeaderFields)-1])
		removed := c.HT.Remove(len(c.HT.HeaderFields) - 1)

		// 消したエントリへの参照を RS からも消す
		c.RS.Remove(removed)
	}
	return
}

// Push new enctory to Header Table
// and Eviction
func (c *Context) Push(hf *HeaderField) {
	c.HT.Push(hf)
	c.Eviction()
}

// Dump for Debug
func (c *Context) Dump() string {
	return fmt.Sprintf("%v%v%v", c.HT.Dump(), c.RS.Dump(), c.ES.Dump())
}

/*
func (c *Context) Encode(header http.Header) []byte {
	var buf bytes.Buffer

	// http.Header を HeaderSet に変換
	headerSet := HeaderToHeaderSet(header)

	// ReferenceSet の中から消すべき値を消す
	buf.Write(c.CleanReferenceSet(headerSet))

	// Header Set の中から送らない値を消す
	c.CleanHeaderSet(headerSet)

	// Header Table にあるやつを処理
	buf.Write(c.ProcessHeader(headerSet))
	return buf.Bytes()
}

// 1. 不要なエントリを reference set から消す
func (c *Context) CleanReferenceSet(headerSet HeaderSet) []byte {
	var buf bytes.Buffer
	// reference set の中にあって、 header set の中に無いものは
	// 相手の reference set から消さないといけないので、
	// indexed representation でエンコードして
	// reference set からは消す
	for name, value := range c.ReferenceSet {
		if headerSet[name] != value {
			c.ReferenceSet.Del(name)

			// Header Table を探して、 index だけ取り出す
			index, _ := c.HeaderTable.SearchHeader(name, value)

			// Indexed Header を生成
			frame := NewIndexedHeader(uint64(index))
			f := frame.Encode()
			buf.Write(f.Bytes())
			Debug(fmt.Sprintf("indexed header index=%v removal from reference set", index))
		}
	}
	return buf.Bytes()
}

// 2. 送る必要の無いものを header set から消す
func (c *Context) CleanHeaderSet(headerSet HeaderSet) {
	for name, value := range c.ReferenceSet {
		if headerSet[name] == value {
			delete(headerSet, name)
			// TODO: "common-header" としてマーク
			Debug(fmt.Sprintf("remove from header set %v %v", name, value))
		}
	}
}

// 3 と 4. 残りの処理
func (c *Context) ProcessHeader(headerSet HeaderSet) []byte {
	var buf bytes.Buffer
	for name, value := range headerSet {
		index, h := c.HeaderTable.SearchHeader(name, value)
		if h != nil { // 3.1 HT にエントリがある
			// Indexed Heaer で index だけ送れば良い
			frame := NewIndexedHeader(uint64(index))
			f := frame.Encode()
			Debug(fmt.Sprintf("indexed header index=%v", index))
			Debug(fmt.Sprintf("add to refset (%q, %q)", name, value))
			c.ReferenceSet.Add(name, value)
			buf.Write(f.Bytes())
		} else if index != -1 { // HT に name だけある
			// Indexed Name Without Indexing
			// value だけ送る。 HT は更新しない。
			frame := NewIndexedNameWithoutIndexing(uint64(index), value)
			f := frame.Encode()
			Debug(fmt.Sprintf("literal header without indexing, name index=%v value=%q", index, value))
			buf.Write(f.Bytes())
		} else { // HT に name も value もない
			// New Name Without Indexing
			// name, value を送って HT は更新しない。
			frame := NewNewNameWithoutIndexing(name, value)
			f := frame.Encode()
			Debug(fmt.Sprintf("literal header without indexing, new name name=%q value=%q", name, value))
			buf.Write(f.Bytes())
		}
	}
	return buf.Bytes()
}
*/
