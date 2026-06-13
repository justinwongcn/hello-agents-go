package main

import (
	"fmt"
	"math"
	"math/rand/v2"
)

// MultiHeadAttention 多头注意力机制模块
type MultiHeadAttention struct {
	dModel   int
	numHeads int
	dK       int
	WQ       [][]float64
	WK       [][]float64
	WV       [][]float64
	WO       [][]float64
}

// NewMultiHeadAttention 创建多头注意力模块
func NewMultiHeadAttention(dModel, numHeads int) *MultiHeadAttention {
	if dModel%numHeads != 0 {
		panic("d_model 必须能被 num_heads 整除")
	}
	dK := dModel / numHeads
	return &MultiHeadAttention{
		dModel:   dModel,
		numHeads: numHeads,
		dK:       dK,
		WQ:       randomMatrix(dModel, dModel),
		WK:       randomMatrix(dModel, dModel),
		WV:       randomMatrix(dModel, dModel),
		WO:       randomMatrix(dModel, dModel),
	}
}

// ScaledDotProductAttention 计算缩放点积注意力
func (m *MultiHeadAttention) ScaledDotProductAttention(Q, K, V [][]float64, mask [][]bool) [][]float64 {
	seqLen := len(Q)
	dK := len(Q[0])

	// 1. 计算注意力得分 (QK^T)
	scores := make([][]float64, seqLen)
	for i := range seqLen {
		scores[i] = make([]float64, seqLen)
		for j := range seqLen {
			dot := 0.0
			for k := range dK {
				dot += Q[i][k] * K[j][k]
			}
			scores[i][j] = dot / math.Sqrt(float64(dK))
		}
	}

	// 2. 应用掩码 (如果提供)
	if mask != nil {
		for i := range seqLen {
			for j := range seqLen {
				if !mask[i][j] {
					scores[i][j] = -1e9
				}
			}
		}
	}

	// 3. 计算注意力权重 (Softmax)
	for i := range seqLen {
		maxVal := scores[i][0]
		for j := 1; j < seqLen; j++ {
			if scores[i][j] > maxVal {
				maxVal = scores[i][j]
			}
		}
		sumExp := 0.0
		for j := range seqLen {
			scores[i][j] = math.Exp(scores[i][j] - maxVal)
			sumExp += scores[i][j]
		}
		for j := range seqLen {
			scores[i][j] /= sumExp
		}
	}

	// 4. 加权求和 (权重 * V)
	output := make([][]float64, seqLen)
	for i := range seqLen {
		output[i] = make([]float64, len(V[0]))
		for j := range len(V[0]) {
			for k := range seqLen {
				output[i][j] += scores[i][k] * V[k][j]
			}
		}
	}
	return output
}

// PositionWiseFeedForward 位置前馈网络模块
type PositionWiseFeedForward struct {
	dModel int
	dFF    int
	dropout float64
	W1     [][]float64
	B1     []float64
	W2     [][]float64
	B2     []float64
}

// NewPositionWiseFeedForward 创建位置前馈网络
func NewPositionWiseFeedForward(dModel, dFF int, dropout float64) *PositionWiseFeedForward {
	return &PositionWiseFeedForward{
		dModel: dModel,
		dFF:    dFF,
		dropout: dropout,
		W1:     randomMatrix(dModel, dFF),
		B1:     make([]float64, dFF),
		W2:     randomMatrix(dFF, dModel),
		B2:     make([]float64, dModel),
	}
}

// Forward 前向传播
// x 形状: (seq_len, d_model)
func (f *PositionWiseFeedForward) Forward(x [][]float64) [][]float64 {
	seqLen := len(x)
	// x 形状: (seq_len, d_model)

	// 线性变换1: x -> (seq_len, d_ff)
	hidden := make([][]float64, seqLen)
	for i := range seqLen {
		hidden[i] = make([]float64, f.dFF)
		for j := range f.dFF {
			sum := f.B1[j]
			for k := range f.dModel {
				sum += x[i][k] * f.W1[k][j]
			}
			// ReLU激活
			if sum < 0 {
				sum = 0
			}
			hidden[i][j] = sum
		}
	}

	// 线性变换2: (seq_len, d_ff) -> (seq_len, d_model)
	output := make([][]float64, seqLen)
	for i := range seqLen {
		output[i] = make([]float64, f.dModel)
		for j := range f.dModel {
			sum := f.B2[j]
			for k := range f.dFF {
				sum += hidden[i][k] * f.W2[k][j]
			}
			output[i][j] = sum
		}
	}
	// 最终输出形状: (seq_len, d_model)
	return output
}

// PositionalEncoding 为输入序列的词嵌入向量添加位置编码
type PositionalEncoding struct {
	dModel int
	pe     [][]float64
}

// NewPositionalEncoding 创建位置编码
func NewPositionalEncoding(dModel, maxLen int) *PositionalEncoding {
	pe := make([][]float64, maxLen)
	for i := range maxLen {
		pe[i] = make([]float64, dModel)
		for j := range dModel {
			if j%2 == 0 {
				// 偶数维度使用 sin
				pe[i][j] = math.Sin(float64(i) * math.Exp(float64(j)*(-math.Log(10000.0)/float64(dModel))))
			} else {
				// 奇数维度使用 cos
				pe[i][j] = math.Cos(float64(i) * math.Exp(float64(j-1)*(-math.Log(10000.0)/float64(dModel))))
			}
		}
	}
	return &PositionalEncoding{dModel: dModel, pe: pe}
}

// Forward 前向传播
func (pe *PositionalEncoding) Forward(x [][]float64) [][]float64 {
	seqLen := len(x)
	output := make([][]float64, seqLen)
	for i := range seqLen {
		output[i] = make([]float64, pe.dModel)
		for j := range pe.dModel {
			output[i][j] = x[i][j] + pe.pe[i][j]
		}
	}
	return output
}

// EncoderLayer 编码器核心层
type EncoderLayer struct {
	selfAttn    *MultiHeadAttention
	feedForward *PositionWiseFeedForward
	norm1       *LayerNorm
	norm2       *LayerNorm
	dropout     float64
}

// NewEncoderLayer 创建编码器层
func NewEncoderLayer(dModel, numHeads, dFF int, dropout float64) *EncoderLayer {
	return &EncoderLayer{
		selfAttn:    NewMultiHeadAttention(dModel, numHeads),
		feedForward: NewPositionWiseFeedForward(dModel, dFF, dropout),
		norm1:       NewLayerNorm(dModel),
		norm2:       NewLayerNorm(dModel),
		dropout:     dropout,
	}
}

// Forward 前向传播
func (el *EncoderLayer) Forward(x [][]float64, mask [][]bool) [][]float64 {
	// 1. 多头自注意力
	attnOutput := el.selfAttn.Forward(x, x, x, mask)
	// 残差连接 + 层归一化
	residual := addMatrices(x, dropout(attnOutput, el.dropout))
	x = el.norm1.Forward(residual)

	// 2. 前馈网络
	ffOutput := el.feedForward.Forward(x)
	residual = addMatrices(x, dropout(ffOutput, el.dropout))
	x = el.norm2.Forward(residual)

	return x
}

// DecoderLayer 解码器核心层
type DecoderLayer struct {
	selfAttn    *MultiHeadAttention
	crossAttn   *MultiHeadAttention
	feedForward *PositionWiseFeedForward
	norm1       *LayerNorm
	norm2       *LayerNorm
	norm3       *LayerNorm
	dropout     float64
}

// NewDecoderLayer 创建解码器层
func NewDecoderLayer(dModel, numHeads, dFF int, dropout float64) *DecoderLayer {
	return &DecoderLayer{
		selfAttn:    NewMultiHeadAttention(dModel, numHeads),
		crossAttn:   NewMultiHeadAttention(dModel, numHeads),
		feedForward: NewPositionWiseFeedForward(dModel, dFF, dropout),
		norm1:       NewLayerNorm(dModel),
		norm2:       NewLayerNorm(dModel),
		norm3:       NewLayerNorm(dModel),
		dropout:     dropout,
	}
}

// Forward 前向传播
func (dl *DecoderLayer) Forward(x, encoderOutput [][]float64, srcMask, tgtMask [][]bool) [][]float64 {
	// 1. 掩码多头自注意力 (对自己)
	attnOutput := dl.selfAttn.Forward(x, x, x, tgtMask)
	residual := addMatrices(x, dropout(attnOutput, dl.dropout))
	x = dl.norm1.Forward(residual)

	// 2. 交叉注意力 (对编码器输出)
	crossAttnOutput := dl.crossAttn.Forward(x, encoderOutput, encoderOutput, srcMask)
	residual = addMatrices(x, dropout(crossAttnOutput, dl.dropout))
	x = dl.norm2.Forward(residual)

	// 3. 前馈网络
	ffOutput := dl.feedForward.Forward(x)
	residual = addMatrices(x, dropout(ffOutput, dl.dropout))
	x = dl.norm3.Forward(residual)

	return x
}

// Encoder 编码器
type Encoder struct {
	embedding  [][]float64
	posEncoder *PositionalEncoding
	layers     []*EncoderLayer
	norm       *LayerNorm
}

// NewEncoder 创建编码器
func NewEncoder(vocabSize, dModel, numLayers, numHeads, dFF int, dropout float64, maxLen int) *Encoder {
	layers := make([]*EncoderLayer, numLayers)
	for i := range numLayers {
		layers[i] = NewEncoderLayer(dModel, numHeads, dFF, dropout)
	}
	return &Encoder{
		embedding:  randomMatrix(vocabSize, dModel),
		posEncoder: NewPositionalEncoding(dModel, maxLen),
		layers:     layers,
		norm:       NewLayerNorm(dModel),
	}
}

// Forward 前向传播
func (e *Encoder) Forward(x [][]int, mask [][]bool) [][]float64 {
	seqLen := len(x[0])
	batchSize := len(x)
	dModel := len(e.embedding[0])

	// 嵌入
	embedded := make([][][]float64, batchSize)
	for b := range batchSize {
		embedded[b] = make([][]float64, seqLen)
		for s := range seqLen {
			embedded[b][s] = make([]float64, dModel)
			copy(embedded[b][s], e.embedding[x[b][s]])
		}
	}

	// 位置编码
	for b := range batchSize {
		embedded[b] = e.posEncoder.Forward(embedded[b])
	}

	// 通过编码器层
	output := embedded[0] // 简化：只处理batch_size=1
	for _, layer := range e.layers {
		output = layer.Forward(output, mask)
	}
	return e.norm.Forward(output)
}

// Decoder 解码器
type Decoder struct {
	embedding  [][]float64
	posEncoder *PositionalEncoding
	layers     []*DecoderLayer
	norm       *LayerNorm
}

// NewDecoder 创建解码器
func NewDecoder(vocabSize, dModel, numLayers, numHeads, dFF int, dropout float64, maxLen int) *Decoder {
	layers := make([]*DecoderLayer, numLayers)
	for i := range numLayers {
		layers[i] = NewDecoderLayer(dModel, numHeads, dFF, dropout)
	}
	return &Decoder{
		embedding:  randomMatrix(vocabSize, dModel),
		posEncoder: NewPositionalEncoding(dModel, maxLen),
		layers:     layers,
		norm:       NewLayerNorm(dModel),
	}
}

// Forward 前向传播
func (d *Decoder) Forward(x [][]int, encoderOutput [][]float64, srcMask, tgtMask [][]bool) [][]float64 {
	seqLen := len(x[0])
	batchSize := len(x)
	dModel := len(d.embedding[0])

	// 嵌入
	embedded := make([][][]float64, batchSize)
	for b := range batchSize {
		embedded[b] = make([][]float64, seqLen)
		for s := range seqLen {
			embedded[b][s] = make([]float64, dModel)
			copy(embedded[b][s], d.embedding[x[b][s]])
		}
	}

	// 位置编码
	for b := range batchSize {
		embedded[b] = d.posEncoder.Forward(embedded[b])
	}

	// 通过解码器层
	output := embedded[0] // 简化：只处理batch_size=1
	for _, layer := range d.layers {
		output = layer.Forward(output, encoderOutput, srcMask, tgtMask)
	}
	return d.norm.Forward(output)
}

// Transformer Transformer模型
type Transformer struct {
	encoder      *Encoder
	decoder      *Decoder
	finalLinear  [][]float64
	tgtVocabSize int
}

// NewTransformer 创建Transformer模型
func NewTransformer(srcVocabSize, tgtVocabSize, dModel, numLayers, numHeads, dFF int, dropout float64, maxLen int) *Transformer {
	return &Transformer{
		encoder:      NewEncoder(srcVocabSize, dModel, numLayers, numHeads, dFF, dropout, maxLen),
		decoder:      NewDecoder(tgtVocabSize, dModel, numLayers, numHeads, dFF, dropout, maxLen),
		finalLinear:  randomMatrix(dModel, tgtVocabSize),
		tgtVocabSize: tgtVocabSize,
	}
}

// GenerateMask 生成掩码
func (t *Transformer) GenerateMask(src, tgt [][]int) ([][]bool, [][]bool) {
	srcLen := len(src[0])
	tgtLen := len(tgt[0])

	// src_mask: (src_len, src_len) - 简化版
	srcMask := make([][]bool, srcLen)
	for i := range srcLen {
		srcMask[i] = make([]bool, srcLen)
		for j := range srcLen {
			srcMask[i][j] = true // 简化：所有位置都可见
		}
	}

	// tgt_mask: (tgt_len, tgt_len) - 下三角矩阵
	tgtMask := make([][]bool, tgtLen)
	for i := range tgtLen {
		tgtMask[i] = make([]bool, tgtLen)
		for j := range tgtLen {
			tgtMask[i][j] = j <= i // 下三角矩阵
		}
	}

	return srcMask, tgtMask
}

// Forward 前向传播
func (t *Transformer) Forward(src, tgt [][]int) [][]float64 {
	srcMask, tgtMask := t.GenerateMask(src, tgt)

	encoderOutput := t.encoder.Forward(src, srcMask)
	decoderOutput := t.decoder.Forward(tgt, encoderOutput, srcMask, tgtMask)

	// 最终线性变换
	seqLen := len(decoderOutput)
	output := make([][]float64, seqLen)
	for i := range seqLen {
		output[i] = make([]float64, t.tgtVocabSize)
		for j := range t.tgtVocabSize {
			sum := 0.0
			for k := range len(decoderOutput[0]) {
				sum += decoderOutput[i][k] * t.finalLinear[k][j]
			}
			output[i][j] = sum
		}
	}
	return output
}

// LayerNorm 层归一化
type LayerNorm struct {
	dModel int
	gamma  []float64
	beta   []float64
	eps    float64
}

// NewLayerNorm 创建层归一化
func NewLayerNorm(dModel int) *LayerNorm {
	gamma := make([]float64, dModel)
	beta := make([]float64, dModel)
	for i := range dModel {
		gamma[i] = 1.0
		beta[i] = 0.0
	}
	return &LayerNorm{
		dModel: dModel,
		gamma:  gamma,
		beta:   beta,
		eps:    1e-5,
	}
}

// Forward 前向传播
func (ln *LayerNorm) Forward(x [][]float64) [][]float64 {
	seqLen := len(x)
	output := make([][]float64, seqLen)
	for i := range seqLen {
		// 计算均值
		mean := 0.0
		for _, v := range x[i] {
			mean += v
		}
		mean /= float64(ln.dModel)

		// 计算方差
		variance := 0.0
		for _, v := range x[i] {
			variance += (v - mean) * (v - mean)
		}
		variance /= float64(ln.dModel)

		// 归一化
		output[i] = make([]float64, ln.dModel)
		for j := range ln.dModel {
			output[i][j] = ln.gamma[j]*(x[i][j]-mean)/math.Sqrt(variance+ln.eps) + ln.beta[j]
		}
	}
	return output
}

// 辅助函数

// randomMatrix 创建随机矩阵
func randomMatrix(rows, cols int) [][]float64 {
	matrix := make([][]float64, rows)
	for i := range rows {
		matrix[i] = make([]float64, cols)
		for j := range cols {
			matrix[i][j] = rand.NormFloat64() * 0.01
		}
	}
	return matrix
}

// addMatrices 矩阵相加
func addMatrices(a, b [][]float64) [][]float64 {
	rows := len(a)
	cols := len(a[0])
	result := make([][]float64, rows)
	for i := range rows {
		result[i] = make([]float64, cols)
		for j := range cols {
			result[i][j] = a[i][j] + b[i][j]
		}
	}
	return result
}

// dropout Dropout操作
func dropout(x [][]float64, rate float64) [][]float64 {
	if rate == 0 {
		return x
	}
	rows := len(x)
	cols := len(x[0])
	result := make([][]float64, rows)
	for i := range rows {
		result[i] = make([]float64, cols)
		for j := range cols {
			if rand.Float64() > rate {
				result[i][j] = x[i][j] / (1 - rate)
			} else {
				result[i][j] = 0
			}
		}
	}
	return result
}

// Forward 对MultiHeadAttention的Forward方法
func (m *MultiHeadAttention) Forward(Q, K, V [][]float64, mask [][]bool) [][]float64 {
	// 线性变换
	Q = matMul(Q, m.WQ)
	K = matMul(K, m.WK)
	V = matMul(V, m.WV)

	// 计算注意力
	attnOutput := m.ScaledDotProductAttention(Q, K, V, mask)

	// 输出投影
	return matMul(attnOutput, m.WO)
}

// matMul 矩阵乘法
func matMul(a, b [][]float64) [][]float64 {
	rows := len(a)
	cols := len(b[0])
	inner := len(b)
	result := make([][]float64, rows)
	for i := range rows {
		result[i] = make([]float64, cols)
		for j := range cols {
			sum := 0.0
			for k := range inner {
				sum += a[i][k] * b[k][j]
			}
			result[i][j] = sum
		}
	}
	return result
}

func main() {
	// --- 演示如何使用模型 ---
	// 1. 定义超参数
	srcVocabSize := 5000
	tgtVocabSize := 5000
	dModel := 512
	numLayers := 6
	numHeads := 8
	dFF := 2048
	dropout := 0.1
	maxLen := 100

	// 2. 实例化模型
	model := NewTransformer(srcVocabSize, tgtVocabSize, dModel, numLayers, numHeads, dFF, dropout, maxLen)

	// 3. 创建模拟输入数据
	// 假设 batch_size=2, src_seq_len=10, tgt_seq_len=12
	src := make([][]int, 2)
	for i := range 2 {
		src[i] = make([]int, 10)
		for j := range 10 {
			src[i][j] = rand.IntN(srcVocabSize) + 1
		}
	}

	tgt := make([][]int, 2)
	for i := range 2 {
		tgt[i] = make([]int, 12)
		for j := range 12 {
			tgt[i][j] = rand.IntN(tgtVocabSize) + 1
		}
	}

	// 4. 模型前向传播
	output := model.Forward(src, tgt)

	// 5. 打印输出形状
	fmt.Printf("模型输出的形状: [%d, %d, %d]\n", 2, len(output), len(output[0]))
	// 预期输出: [2, 12, 5000] -> (batch_size, tgt_seq_len, tgt_vocab_size)
}
