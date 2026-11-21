# Factorio AI Engineer

## Goal
Claude CodeからFactorio LLM Use MCPで工場を調査・設計・建設する。

## 工場設計の原則
工場の新設・増設・物流変更を提案するときや、次に取り組む課題を決めるときは、先に `docs/philosophy.md` を読む。提案では、**何を実現したいか、何がそれを妨げていると考えるか、どう確かめるか**を簡潔に示す。適用する原則は、その判断に関係するものを選ぶ。

## Environment
- macOS M3 / Factorio 2.1.24 / AngelBob + Space Age
- MCP: factorio-llm-use（RCON経由、セーブのコピーで動作）
- Helmod/YARMは人間用。Claudeは参照しない
- 詳細: docs/setup.md, docs/mods.md（必要時のみ）

## Hard Rules
- レシピは list_recipes の値のみ使用。取得不可なら停止して報告。
- list_recipes は具体的な製品名で query を絞る（"angels" 等の広い語は禁止）。
- 全ツール呼び出しで surface を明示（省略時nauvis）。
- 建設は create_plan → preview_plan → 承認 → apply_plan(ghosts) のみ。
- execute_lua は読み取り専用の調査に限定。状態変更は明示許可時のみ。
- 撤去（mine_entity等）は件数に関係なく承認必須。
- observe は area/radius と names/types で絞る。detailed=true は必要時のみ。

## Workflow
1. 調査: get_world / observe（範囲限定）/ list_recipes
2. 計画: plans/<name>.md にテンプレで記載 → create_plan
3. 確認: preview_plan → 承認待ち
4. 建設: apply_plan(mode="ghosts")
5. 検証: get_action_status / inspect_entity で稼働・詰まりを確認し目標値(/分)と比較

## plans/ テンプレ
- 目的 / surface / 範囲(座標) / 目標生産量
- 使用レシピ（list_recipes値）/ 必要建物数・資材
- 撤去対象 / 検証方法

## Output
- 鉱床・設備・建設候補地などの場所を示すときは、地点名の下に `[gps=x,y,surface]` 形式のタグを、タグだけのコードブロックで出す。座標とサーフェスにはゲーム内で確認した値を使う。
- GitHub Issueの本文・コメント・進捗報告では、冒頭の1〜2文に結論を簡潔に書く。そのあとに、必要な根拠や詳細を続ける。
- ゲーム内の作業についてIssueにコメントする前に、`docs/philosophy.md` を読む。冒頭の1〜2文で結論を述べ、その直後に **原典の適用** として、関連する章の見出しと、その原則を今回の判断や確認にどう使ったかを1文で書く。
成功時は結果のみ。失敗時は原因と次の対応。

## ゲーム内の進捗記録

Factorioの現況調査、設計、建設、動作確認を行ったら、毎回 `PROGRESS.md` の末尾に英語で要約を追記する。日付と対象を見出しにし、結果と、必要なら主な詰まりや次の一手を1〜3文にまとめる。進展がなければ、その理由を一言で残す。調査データ、座標、レシピ表、作業手順などの詳細は、関連するIssue・PR・設計文書へのリンクで参照できるようにする。
