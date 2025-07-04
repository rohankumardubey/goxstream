<h1 class="code-line" data-line-start=0 data-line-end=1 ><a id="GoXStream_0"></a>GoXStream</h1>
<p class="has-line-data" data-line-start="2" data-line-end="5"><strong>GoXStream</strong> is a modern, extensible real-time streaming engine built in Go, inspired by Apache Flink.<br>
It enables processing and transforming streaming data with powerful, configurable pipelines using operators like map, filter, reduce, and both tumbling and sliding windows.<br>
GoXStream supports dynamic job definitions via REST API and is ready for integration with a React UI dashboard.</p>
<hr>
<h2 class="code-line" data-line-start=8 data-line-end=9 ><a id="_Features_8"></a>✨ Features</h2>
<h2 class="code-line" data-line-start=10 data-line-end=11 ><a id="Features_10"></a>Features</h2>
<ul>
<li class="has-line-data" data-line-start="12" data-line-end="13"><strong>Dynamic pipeline definition via REST API</strong></li>
<li class="has-line-data" data-line-start="13" data-line-end="14"><strong>Event-time tumbling &amp; sliding windows</strong></li>
<li class="has-line-data" data-line-start="14" data-line-end="15"><strong>Watermark/late event support</strong> (true stream semantics)</li>
<li class="has-line-data" data-line-start="15" data-line-end="16"><strong>Chaining map, filter, reduce operators</strong></li>
<li class="has-line-data" data-line-start="16" data-line-end="17"><strong>CSV file source/sink (DB/Kafka coming soon)</strong></li>
<li class="has-line-data" data-line-start="17" data-line-end="18"><strong>Ready for React UI dashboard integration</strong></li>
<li class="has-line-data" data-line-start="18" data-line-end="20"><strong>Easy extension: Add custom operators and connectors</strong></li>
</ul>
<hr>
<h2 class="code-line" data-line-start=22 data-line-end=23 ><a id="_Quick_Start_22"></a>🚀 Quick Start</h2>
<h3 class="code-line" data-line-start=24 data-line-end=25 ><a id="1_Clone_and_Build_24"></a>1. <strong>Clone and Build</strong></h3>
<pre><code class="has-line-data" data-line-start="27" data-line-end="31" class="language-bash">git <span class="hljs-built_in">clone</span> https://github.com/YOUR_GITHUB_USERNAME/goxstream.git
<span class="hljs-built_in">cd</span> goxstream
go mod tidy
</code></pre>
<h3 class="code-line" data-line-start=31 data-line-end=32 ><a id="2_Prepare_Input_Data_31"></a>2. Prepare Input Data</h3>
<p class="has-line-data" data-line-start="32" data-line-end="33"><strong><em>Place an example input.csv in the project root:</em></strong></p>
<pre><code class="has-line-data" data-line-start="35" data-line-end="44" class="language-bash">id,name,city
<span class="hljs-number">1</span>,Alice,London
<span class="hljs-number">2</span>,Bob,Berlin
<span class="hljs-number">3</span>,Charlie,Paris
<span class="hljs-number">4</span>,David,Berlin
<span class="hljs-number">5</span>,Eve,Paris
<span class="hljs-number">6</span>,Frank,Paris
<span class="hljs-number">7</span>,Grace,Berlin
</code></pre>
<h3 class="code-line" data-line-start=45 data-line-end=46 ><a id="3_Run_the_API_Server_45"></a>3. Run the API Server</h3>
<pre><code class="has-line-data" data-line-start="47" data-line-end="49" class="language-bash">go run ./cmd/goxstream/main.go
</code></pre>
<h3 class="code-line" data-line-start=50 data-line-end=51 ><a id="4_Submit_a_Pipeline_Job_50"></a>4. Submit a Pipeline Job</h3>
<p class="has-line-data" data-line-start="51" data-line-end="53"><strong><em>Use curl or Postman to submit a dynamic pipeline (example: sliding window reduce):</em></strong><br>
<strong><em>a. Regular time window:</em></strong></p>
<pre><code class="has-line-data" data-line-start="56" data-line-end="75" class="language-bash">curl -X POST http://localhost:<span class="hljs-number">8080</span>/<span class="hljs-built_in">jobs</span> \
  -H <span class="hljs-string">"Content-Type: application/json"</span> \
  <span class="hljs-operator">-d</span> <span class="hljs-string">'{
    "source": { "type": "file", "path": "input.csv" },
    "operators": [
      {
        "type": "time_window",
        "params": {
          "duration": "10s",
          "inner": {
            "type": "reduce",
            "params": { "key": "city", "agg": "count" }
          }
        }
      }
    ],
    "sink": { "type": "file", "path": "output.csv" }
  }'</span>
</code></pre>
<p class="has-line-data" data-line-start="76" data-line-end="77"><strong><em>b. Time window with watermark/late event support:</em></strong></p>
<pre><code class="has-line-data" data-line-start="78" data-line-end="99" class="language-bash">curl -X POST http://localhost:<span class="hljs-number">8080</span>/<span class="hljs-built_in">jobs</span> \
  -H <span class="hljs-string">"Content-Type: application/json"</span> \
  <span class="hljs-operator">-d</span> <span class="hljs-string">'{
    "source": { "type": "file", "path": "input.csv" },
    "operators": [
      {
        "type": "time_window_watermark",
        "params": {
          "duration": "10s",
          "allowed_lateness": "5s",
          "inner": {
            "type": "reduce",
            "params": { "key": "city", "agg": "count" }
          }
        }
      }
    ],
    "sink": { "type": "file", "path": "output.csv" }
  }'</span>

</code></pre>
<p class="has-line-data" data-line-start="99" data-line-end="100"><strong><em>Your results will be in output.csv with a window_id column.</em></strong></p>
<hr>
<h3 class="code-line" data-line-start=103 data-line-end=104 ><a id="_Architecture_103"></a>🛠️ Architecture</h3>
<pre><code class="has-line-data" data-line-start="105" data-line-end="110" class="language-bash">[Source] --&gt; [Map] --&gt; [Filter] --&gt; [Window/Reduce] --&gt; [Sink]
   |            |         |           |                   |
 [File]   [Add/Transform] [Select] [Tumbling/Sliding]  [File]

</code></pre>
<p class="has-line-data" data-line-start="111" data-line-end="112"><strong><em>Operators: Implemented as Go interfaces, dynamically composed at runtime.</em></strong></p>
<p class="has-line-data" data-line-start="113" data-line-end="114"><strong><em>REST API: Accepts JSON job specs, launches pipeline as background Go routines.</em></strong></p>
<hr>
<h3 class="code-line" data-line-start=116 data-line-end=117 ><a id="_JSON_Job_Spec_116"></a>📝 JSON Job Spec</h3>
<p class="has-line-data" data-line-start="117" data-line-end="118"><strong><em>A pipeline is defined by a simple JSON:</em></strong></p>
<pre><code class="has-line-data" data-line-start="120" data-line-end="141" class="language-bash">{
  <span class="hljs-string">"source"</span>: {<span class="hljs-string">"type"</span>: <span class="hljs-string">"file"</span>, <span class="hljs-string">"path"</span>: <span class="hljs-string">"input.csv"</span>},
  <span class="hljs-string">"operators"</span>: [
    {<span class="hljs-string">"type"</span>: <span class="hljs-string">"map"</span>, <span class="hljs-string">"params"</span>: {<span class="hljs-string">"col"</span>: <span class="hljs-string">"processed"</span>, <span class="hljs-string">"val"</span>: <span class="hljs-string">"yes"</span>}},
    {<span class="hljs-string">"type"</span>: <span class="hljs-string">"filter"</span>, <span class="hljs-string">"params"</span>: {<span class="hljs-string">"field"</span>: <span class="hljs-string">"city"</span>, <span class="hljs-string">"eq"</span>: <span class="hljs-string">"Berlin"</span>}},
    {
      <span class="hljs-string">"type"</span>: <span class="hljs-string">"sliding_window"</span>,
      <span class="hljs-string">"params"</span>: {
        <span class="hljs-string">"size"</span>: <span class="hljs-number">3</span>,
        <span class="hljs-string">"step"</span>: <span class="hljs-number">1</span>,
        <span class="hljs-string">"inner"</span>: {
          <span class="hljs-string">"type"</span>: <span class="hljs-string">"reduce"</span>,
          <span class="hljs-string">"params"</span>: {<span class="hljs-string">"key"</span>: <span class="hljs-string">"city"</span>, <span class="hljs-string">"agg"</span>: <span class="hljs-string">"count"</span>}
        }
      }
    }
  ],
  <span class="hljs-string">"sink"</span>: {<span class="hljs-string">"type"</span>: <span class="hljs-string">"file"</span>, <span class="hljs-string">"path"</span>: <span class="hljs-string">"output.csv"</span>}
}

</code></pre>
<hr>
<h3 class="code-line" data-line-start=145 data-line-end=146 ><a id="_Operator_Types_145"></a>📚 Operator Types</h3>
<pre><code class="has-line-data" data-line-start="148" data-line-end="156" class="language-bash">| Type             | Description              | Example Params                        |
| ---------------- | ------------------------ | ------------------------------------- |
| map              | Add or transform columns | `col`, `val`                          |
| filter           | Filter rows by condition | `field`, `eq`                         |
| reduce           | Aggregate/group by field | `key`, `agg` (`count`, future: `sum`) |
| tumbling\_window | Non-overlapping windows  | `size`, `inner`                       |
| sliding\_window  | Overlapping windows      | `size`, `step`, `inner`               |
</code></pre>
<hr>
<h3 class="code-line" data-line-start=159 data-line-end=160 ><a id="_Extending_GoXStream_159"></a>🧑‍💻 Extending GoXStream</h3>
<p class="has-line-data" data-line-start="160" data-line-end="162"><strong><em>Add New Operators:<br>
Implement the Operator interface and add a factory to the operator registry.</em></strong></p>
<p class="has-line-data" data-line-start="163" data-line-end="165"><strong><em>Support New Sources/Sinks:<br>
Implement a Source or Sink interface in internal/source or internal/sink.</em></strong></p>
<p class="has-line-data" data-line-start="166" data-line-end="168"><strong><em>React UI Integration:<br>
Planned for interactive pipeline creation and monitoring.</em></strong></p>
<hr>
<h3 class="code-line" data-line-start=171 data-line-end=172 ><a id="_Roadmap_171"></a>🔜 Roadmap</h3>
<p class="has-line-data" data-line-start="172" data-line-end="173"><input type="checkbox" id="checkbox21" checked="true"><label for="checkbox21">Count-based tumbling/sliding windows</label></p>
<p class="has-line-data" data-line-start="174" data-line-end="175"><input type="checkbox" id="checkbox22" checked="true"><label for="checkbox22">Time-based tumbling/sliding windows</label></p>
<p class="has-line-data" data-line-start="176" data-line-end="177"><input type="checkbox" id="checkbox23" checked="true"><label for="checkbox23">Watermark &amp; late event support</label></p>
<p class="has-line-data" data-line-start="178" data-line-end="179"><input type="checkbox" id="checkbox24"><label for="checkbox24">DB &amp; Kafka sources/sinks</label></p>
<p class="has-line-data" data-line-start="180" data-line-end="181"><input type="checkbox" id="checkbox25"><label for="checkbox25">React UI dashboard</label></p>
<p class="has-line-data" data-line-start="182" data-line-end="183"><input type="checkbox" id="checkbox26"><label for="checkbox26">More aggregations: sum, avg, min, max</label></p>
<p class="has-line-data" data-line-start="184" data-line-end="185"><input type="checkbox" id="checkbox27"><label for="checkbox27">State, session windows, custom UDFs</label></p>
<hr>
<h3 class="code-line" data-line-start=189 data-line-end=190 ><a id="_Contributing_189"></a>🙌 Contributing</h3>
<p class="has-line-data" data-line-start="190" data-line-end="192"><strong><em>PRs, issues, and ideas are welcome!<br>
Fork and submit improvements or new features—let’s build a great Go stream engine together!</em></strong></p>
<h3 class="code-line" data-line-start=194 data-line-end=195 ><a id="GoXStream__Streaming_the_Go_way__194"></a>GoXStream — Streaming, the Go way! 🚀</h3>