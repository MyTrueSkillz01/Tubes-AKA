package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"runtime"
	"time"
)

type TuringMachine struct {
	Tape []byte
	Head int
	Len  int
}

func NewTM_Benchmark(size int) *TuringMachine {
	tape := make([]byte, size)
	pattern := []byte("IFMMP")
	patLen := len(pattern)
	for i := 0; i < size; i++ {
		tape[i] = pattern[i%patLen]
	}
	return &TuringMachine{Tape: tape, Head: 0, Len: size}
}

func NewTM_Manual(input string) *TuringMachine {
	tape := []byte(input)
	return &TuringMachine{Tape: tape, Head: 0, Len: len(tape)}
}

func transition(b byte) byte {
	if b >= 'A' && b <= 'Z' {
		if b == 'A' {
			return 'Z'
		}
		return b - 1
	}
	if b >= 'a' && b <= 'z' {
		if b == 'a' {
			return 'z'
		}
		return b - 1
	}
	return b
}

// --- ITERATIF ---
func (tm *TuringMachine) RunIterative() {
	for tm.Head < tm.Len {
		tm.Tape[tm.Head] = transition(tm.Tape[tm.Head])
		tm.Head++
	}
}

// --- REKURSIF ---
func (tm *TuringMachine) recursiveStep(head int) {
	if head >= tm.Len {
		return
	}
	tm.Tape[head] = transition(tm.Tape[head])
	tm.recursiveStep(head + 1)
}

func (tm *TuringMachine) RunRecursive() (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("STACK OVERFLOW")
		}
	}()
	tm.recursiveStep(0)
	return nil
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.New("index").Parse(htmlTemplate)
		t.Execute(w, nil)
	})

	http.HandleFunc("/api/benchmark", handleBenchmarkStep)

	http.HandleFunc("/api/manual", handleManualInput)

	fmt.Println("🚀 Aplikasi Siap! Buka: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

type ManualReq struct {
	Text string `json:"text"`
}
type ManualRes struct {
	IterOut  string  `json:"iter_out"`
	IterTime float64 `json:"iter_time"`
	RecOut   string  `json:"rec_out"`
	RecTime  float64 `json:"rec_time"`
	RecError string  `json:"rec_error"`
}

func handleManualInput(w http.ResponseWriter, r *http.Request) {
	var req ManualReq
	json.NewDecoder(r.Body).Decode(&req)

	res := ManualRes{}

	tmIter := NewTM_Manual(req.Text)
	start := time.Now()
	tmIter.RunIterative()
	res.IterTime = float64(time.Since(start).Microseconds()) / 1000.0
	res.IterOut = string(tmIter.Tape)

	if len(req.Text) > 50000 {
		res.RecError = "Input terlalu panjang untuk Rekursif Manual (Gunakan Benchmark)"
	} else {
		tmRec := NewTM_Manual(req.Text)
		startRec := time.Now()
		err := tmRec.RunRecursive()
		if err != nil {
			res.RecError = "Stack Overflow"
		} else {
			res.RecTime = float64(time.Since(startRec).Microseconds()) / 1000.0
			res.RecOut = string(tmRec.Tape)
		}
	}

	json.NewEncoder(w).Encode(res)
}

type BenchRequest struct {
	Size int `json:"size"`
}
type BenchResult struct {
	Size      int     `json:"size"`
	IterTime  float64 `json:"iter_time"`
	RecTime   float64 `json:"rec_time"`
	RecError  string  `json:"rec_error"`
	Formatted string  `json:"formatted"`
}

func handleBenchmarkStep(w http.ResponseWriter, r *http.Request) {
	var req BenchRequest
	json.NewDecoder(r.Body).Decode(&req)
	size := req.Size
	result := BenchResult{Size: size, Formatted: formatNumber(size)}

	runtime.GC()

	tmIter := NewTM_Benchmark(size)
	start := time.Now()
	tmIter.RunIterative()
	result.IterTime = float64(time.Since(start).Microseconds()) / 1000.0
	tmIter = nil
	runtime.GC()

	LIMIT_AMAN_REKURSI := 15_000_000
	if size > LIMIT_AMAN_REKURSI {
		result.RecTime = -1
		result.RecError = "Skipped (Too Big)"
	} else {
		tmRec := NewTM_Benchmark(size)
		startRec := time.Now()
		err := tmRec.RunRecursive()
		if err != nil {
			result.RecTime = -1
			result.RecError = "Stack Overflow"
		} else {
			result.RecTime = float64(time.Since(startRec).Microseconds()) / 1000.0
		}
	}
	json.NewEncoder(w).Encode(result)
}

func formatNumber(n int) string {
	if n >= 1_000_000 {
		return fmt.Sprintf("%d Juta", n/1_000_000)
	}
	if n >= 1_000 {
		return fmt.Sprintf("%d Ribu", n/1_000)
	}
	return fmt.Sprintf("%d", n)
}

const htmlTemplate = `
<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Turing Machine Simulator</title>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <style>
        :root { --bg: #0d1117; --card: #161b22; --text: #c9d1d9; --accent: #58a6ff; --green: #238636; --red: #da3633; --border: #30363d; }
        body { font-family: 'Segoe UI', sans-serif; background: var(--bg); color: var(--text); padding: 20px; display: flex; flex-direction: column; align-items: center; }
        .container { width: 100%; max-width: 1000px; }
        
        /* HEADER & TABS */
        h1 { color: var(--accent); text-align: center; margin-bottom: 5px; }
        .tabs { display: flex; justify-content: center; margin-bottom: 20px; gap: 10px; }
        .tab-btn { background: var(--card); color: var(--text); border: 1px solid var(--border); padding: 10px 20px; cursor: pointer; border-radius: 6px; font-weight: bold; }
        .tab-btn.active { background: var(--accent); color: white; border-color: var(--accent); }
        
        /* PANELS */
        .panel { display: none; animation: fadeIn 0.3s; }
        .panel.active { display: block; }
        @keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }

        /* CARDS & INPUTS */
        .card { background: var(--card); padding: 20px; border-radius: 8px; border: 1px solid var(--border); margin-bottom: 20px; }
        textarea { width: 100%; background: #0d1117; color: white; border: 1px solid var(--border); padding: 10px; border-radius: 5px; font-family: monospace; resize: vertical; box-sizing: border-box; }
        .btn-action { width: 100%; padding: 12px; margin-top: 10px; background: var(--green); border: none; color: white; font-weight: bold; cursor: pointer; border-radius: 5px; font-size: 1rem; }
        .btn-action:hover { background: #2ea043; }
        
        /* RESULTS GRID */
        .res-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }
        .res-box h3 { margin: 0 0 10px 0; border-bottom: 1px solid var(--border); padding-bottom: 5px; font-size: 1rem; }
        .iter-head { color: var(--green); }
        .rec-head { color: var(--red); }
        .time-tag { float: right; font-size: 0.8em; background: var(--border); padding: 2px 6px; border-radius: 4px; color: white; }

        /* BENCHMARK STYLES */
        .bench-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 20px; }
        table { width: 100%; border-collapse: collapse; font-size: 0.9em; }
        th, td { padding: 8px; text-align: left; border-bottom: 1px solid var(--border); }
        .status-ok { color: var(--green); font-weight: bold; }
        .status-err { color: var(--red); font-weight: bold; }
        .loading { color: #e3b341; }
    </style>
</head>
<body>

<div class="container">
    <h1>🔐 Turing Machine Simulator</h1>
    <p style="text-align: center; color: #8b949e; margin-bottom: 20px;">Caesar Cipher Decryptor (Iteratif vs Rekursif)</p>

    <div class="tabs">
        <button class="tab-btn active" onclick="switchTab('manual')">✍️ Mode Manual</button>
        <button class="tab-btn" onclick="switchTab('bench')">📊 Mode Benchmark</button>
    </div>

    <div id="manual" class="panel active">
        <div class="card">
            <h3 style="margin-top:0">Input Teks Sandi</h3>
            <textarea id="inputText" rows="3" placeholder="Contoh: IFMMP (akan menjadi HELLO) atau ketik kalimat sandi Anda..."></textarea>
            <button class="btn-action" onclick="runManual()">▶ DEKRIPSI SEKARANG</button>
        </div>

        <div class="res-grid">
            <div class="card">
                <h3 class="iter-head">🔵 Iteratif <span id="iterTime" class="time-tag">0 ms</span></h3>
                <textarea id="iterOut" rows="4" readonly placeholder="Hasil Iteratif..."></textarea>
            </div>
            <div class="card">
                <h3 class="rec-head">🔴 Rekursif <span id="recTime" class="time-tag">0 ms</span></h3>
                <textarea id="recOut" rows="4" readonly placeholder="Hasil Rekursif..."></textarea>
            </div>
        </div>
    </div>

    <div id="bench" class="panel">
        <button class="btn-action" id="startBtn" onclick="startBenchmark()" style="margin-bottom: 20px; background: var(--accent);">🚀 Jalankan Stress Test (100 - 100 Juta)</button>
        
        <div class="bench-grid">
            <div class="card">
                <h3 style="margin-top:0">📝 Log Data</h3>
                <table>
                    <thead><tr><th>Input</th><th>Iteratif</th><th>Rekursif</th></tr></thead>
                    <tbody id="logTable"><tr><td colspan="3" style="text-align:center; color:#555">Menunggu...</td></tr></tbody>
                </table>
            </div>
            <div class="card">
                <h3 style="margin-top:0">📈 Grafik</h3>
                <canvas id="myChart"></canvas>
            </div>
        </div>
    </div>
</div>

<script>
    // --- TAB LOGIC ---
    function switchTab(id) {
        document.querySelectorAll('.panel').forEach(p => p.classList.remove('active'));
        document.querySelectorAll('.tab-btn').forEach(b => b.classList.remove('active'));
        document.getElementById(id).classList.add('active');
        event.target.classList.add('active');
    }

    // --- MANUAL LOGIC ---
    async function runManual() {
        const txt = document.getElementById('inputText').value;
        if(!txt) return alert("Masukkan teks terlebih dahulu!");

        const btn = document.querySelector('.btn-action');
        btn.disabled = true; btn.innerText = "Memproses...";

        try {
            const res = await fetch('/api/manual', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify({text: txt})
            });
            const data = await res.json();

            document.getElementById('iterOut').value = data.iter_out;
            document.getElementById('iterTime').innerText = data.iter_time + " ms";

            if(data.rec_error) {
                document.getElementById('recOut').value = "❌ ERROR: " + data.rec_error;
                document.getElementById('recTime').innerText = "GAGAL";
            } else {
                document.getElementById('recOut').value = data.rec_out;
                document.getElementById('recTime').innerText = data.rec_time + " ms";
            }
        } catch(e) {
            alert("Terjadi kesalahan koneksi");
        }
        btn.disabled = false; btn.innerText = "▶ DEKRIPSI SEKARANG";
    }

    // --- BENCHMARK LOGIC ---
    const milestones = [100, 1000, 5000, 10000, 50000, 100000, 500000, 1000000, 5000000, 10000000, 50000000, 100000000];
    let myChart = null;

    function initChart() {
        const ctx = document.getElementById('myChart').getContext('2d');
        if(myChart) myChart.destroy();
        myChart = new Chart(ctx, {
            type: 'line',
            data: { labels: [], datasets: [
                { label: 'Iteratif', borderColor: '#238636', data: [], tension: 0.1 },
                { label: 'Rekursif', borderColor: '#da3633', data: [], tension: 0.1 }
            ]},
            options: { responsive: true, scales: { x: { display: false }, y: { beginAtZero: true } } }
        });
    }

    async function startBenchmark() {
        const btn = document.getElementById('startBtn');
        const tbody = document.getElementById('logTable');
        btn.disabled = true; btn.innerText = "⏳ Sedang Berjalan...";
        tbody.innerHTML = "";
        initChart();

        for (const size of milestones) {
            const row = tbody.insertRow();
            row.innerHTML = "<td>"+size.toLocaleString()+"</td><td class='loading'>...</td><td class='loading'>...</td>";

            const resp = await fetch('/api/benchmark', {
                method: 'POST', body: JSON.stringify({ size: size })
            });
            const data = await resp.json();

            row.cells[1].innerHTML = "<span class='status-ok'>" + data.iter_time.toFixed(3) + " ms</span>";
            if (data.rec_time === -1) {
                row.cells[2].innerHTML = "<span class='status-err'>" + data.rec_error + "</span>";
                myChart.data.datasets[1].data.push(null);
            } else {
                row.cells[2].innerHTML = "<span class='status-ok'>" + data.rec_time.toFixed(3) + " ms</span>";
                myChart.data.datasets[1].data.push(data.rec_time);
            }
            
            myChart.data.labels.push(data.formatted);
            myChart.data.datasets[0].data.push(data.iter_time);
            myChart.update();
        }
        btn.disabled = false; btn.innerText = "✅ Selesai - Ulangi?";
    }
</script>

</body>
</html>
`
