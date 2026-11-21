# Factorio AI Engineer

## Goal
Claude CodeからFactorio LLM Use MCPで工場を調査・設計・建設する。

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
- GitHub Issueの本文・コメント・進捗報告では、冒頭の1〜2文に結論を簡潔に書く。そのあとに、必要な根拠や詳細を続ける。
成功時は結果のみ。失敗時は原因と次の対応。
