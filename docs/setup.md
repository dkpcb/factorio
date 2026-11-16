# Setup

## 起動（ターミナル1、開いたまま）
cd ~/factorio-llm-use
uv run python scripts/start_server.py --resume \
  --factorio "$HOME/Library/Application Support/Steam/steamapps/common/Factorio/factorio.app/Contents/MacOS/factorio"
# 初回のみ: --resume の代わりに --save ".../saves/FOR_AI.zip" --mods ".../factorio/mods"

## 保存（ターミナル2）
cd ~/factorio-llm-use && uv run python scripts/save_session.py
# 保存完了の表示後、ターミナル1で Ctrl+C

## 注意
- MOD zipは scripts/package_mod.py（ディレクトリエントリ追加版）で作る
- 進行は ~/factorio-llm-use/.runtime/data/saves/ に保存。元のFOR_AIは変更されない
- 人間参加: Multiplayer > Connect to address > 127.0.0.1:34199
