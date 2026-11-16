# 序盤の銅板ガイド（スティラタイト → copper-plate）

> **状態: 未検証。** この文書は Mac側MCP（list_recipes / observe）に接続できない環境で作成した。
> レシピ名・数量・副産物は確認できていないため、ここでは断定しない。
> 確認済みの事実は `docs/mods.md` と `plans/iron-early.md` にあるものだけ。

## 確認済みの事実（plans/iron-early.md より）
- 鉄板ルートは `angels-ore1`（サフィライト）→ `angels-ore1-crushed` → `iron-plate`。
- 序盤はバーナー機のみ使える（electric-mining-drill / angels-ore-crusher / steel-furnace は未解禁）。
- 炉は `angels-stone-crushed` を受け付けない。砕石は箱に溜めて、手作業で `stone` に変換する。
- 使える建物: burner-mining-drill, angels-burner-ore-crusher, stone-furnace, burner-inserter, wooden-chest（すべて enabled=true）。

## 未検証（ゲーム内で必ず確認すること）
- スティラタイトの item 名（`angels-ore3` と思われるが未確認）。
- `copper-plate` を直接作るレシピが、粉砕鉱石から炉で作れるのか、別の精製工程（angelsrefining の分級・浮選など）が必要なのか。
- 材料数・出力数・副産物の種類と量。
- 必要な設備と、その設備が序盤の研究で解禁済みかどうか。
- スティラタイトの鉱床が近くにあるか（`observe` で angels-ore3 を確認）。無い場合は銅を含む別鉱石も調べる。

## 調査手順（Mac側で実施）
1. `list_recipes`（surface=nauvis）を具体名で絞って実行する。広い語は使わない。
   - `query="copper-plate"`: 製品側のレシピ
   - `query="angels-ore3"`: スティラタイト（名前が違えば `get_world` / `observe` で実名を確認）
   - 粉砕物なら `query="angels-ore3-crushed"`
2. 各レシピの category / energy / 入力 / 出力 / enabled を表にまとめる。`enabled=false` は序盤では使えない。
3. `observe`（area/radius と names で限定）でスティラタイト鉱床の位置・量を確認する。

## 手作業の最短手順（雛形）
レシピ確認後に埋める。想定される形は鉄板と同じ流れ。
1. burner-mining-drill で鉱石を採掘（または手掘り）。
2. 必要なら angels-burner-ore-crusher で粉砕。
3. stone-furnace で `copper-plate` に精錬（燃料は coal）。
4. 副産物（砕石など）は箱に回収。

## 最小限の自動化（雛形）
鉄板ライン（`plans/iron-early.md`）と同じ構成を流用する案:
ドリル → クラッシャー → burner-inserter → 箱 → burner-inserter → 炉 → burner-inserter → 出力箱。
ドリル・炉の台数は、確認したレシピ値と目標生産量（/分）から計算する。数量はレシピ確認後に決める。

## 副産物が詰まるとき
- 炉が受け付けない副産物は、炉に直接入れず間に木箱（wooden-chest）を挟んで受ける（鉄板ラインと同じ方式）。
- 箱が満杯になる前に手作業で回収・変換する。満杯までの時間は、副産物の量（未検証）で決まる。
- 詰まったら、まず出力スロット・箱の空き・燃料切れを確認する。

## 次のステップ
上の「調査手順」を実行し、結果で本書の未検証欄を埋める。その後 `plans/copper-early.md` を作成する（CLAUDE.md のテンプレに従う）。
