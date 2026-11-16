# 序盤の銅板ガイド（スティラタイト → copper-plate）

> **状態: レシピは確認済み。** list_recipes / observe の値を使っている（nauvis、tick 約553,000 時点）。
> 建物の速度・燃料消費だけは未検証で、鉄板ライン（`plans/iron-early.md`）の建設後に inspect_entity で確認する。

## 結論
鉄板とまったく同じ流れで作れる。精製工程（分級・浮選など）は要らない。

**angels-ore3 → angels-burner-ore-crusher で粉砕 → stone-furnace で精錬 → copper-plate**

## 使用レシピ（list_recipes値、すべて enabled=true）
| レシピ | カテゴリ | energy | 入力 | 出力 |
|---|---|---|---|---|
| angels-ore3-crushed | angels-ore-refining-t1 | 1 | angels-ore3 ×2 | angels-ore3-crushed ×2, angels-stone-crushed ×1 |
| angels-ore3-crushed-hand | hand-crafting | 4 | angels-ore3 ×2 | angels-ore3-crushed ×2, angels-stone-crushed ×1 |
| copper-plate | smelting | 7 | angels-ore3-crushed ×3 | copper-plate ×2 |
| angels-stone-from-crushed-stone | crafting | 0.5 | angels-stone-crushed ×2 | stone ×1 |

- 有効な copper-plate のレシピは上の1つだけ。生の angels-ore3 は炉に入れても精錬できない。必ず粉砕してから炉に入れる。
- 鉱石2個から、粉砕鉱石2個と砕石1個ができる。つまり鉱石3個で銅板2枚と砕石1.5個。
- item 名 `angels-ore3` は list_prototypes で存在を確認した。ただし表示名（スティラタイト）はツールでは確かめられないので、Angel's の命名規則に基づく推定。

## 鉱床（observe値）
- angels-ore3 は、原点から半径160の範囲に 1,182 タイルある。
- 範囲は概ね x-45〜-13, y15〜70。原点から南西へ約20〜70タイルの位置で、鉄の鉱床（東側 x52〜90）とは反対側にある。
- 最も濃い場所はまだ調べていない。`plans/copper-early.md` を作るときに、範囲を絞った observe で探す。

## 手作業の最短手順
1. angels-ore3 を手で掘る。
2. 手作業で `angels-ore3-crushed-hand` を作る（鉱石2個で4秒）。
3. 粉砕鉱石を stone-furnace に入れる（燃料は coal）。粉砕鉱石3個で銅板2枚、7秒かかる。
4. 手元に溜まる砕石は `angels-stone-from-crushed-stone` で stone に変えられる。

手作業での粉砕は遅い（鉱石2個に4秒）。早めに angels-burner-ore-crusher を置く。

## 最小限の自動化
鉄板ライン（`plans/iron-early.md`）と同じ構成を、そのまま angels-ore3 の鉱床に置く。

ドリル×2 → クラッシャー → burner-inserter → 箱 → burner-inserter → 炉 → burner-inserter → 出力箱

| 工程 | 1台の能力 | 15枚/分に必要な台数 |
|---|---|---|
| stone-furnace | 2枚/7秒 = 約17.1枚/分（crafting speed 1 と仮定、未検証） | 1 |
| angels-burner-ore-crusher | 鉱石2個/1秒 = 鉱石120個/分（速度1と仮定、未検証） | 1 |
| burner-mining-drill | 約15個/分（未検証） | 2（鉱石 22.5個/分） |

建物の材料（list_recipes値）: burner-mining-drill = ギア3 + 石の炉1 + 鉄板3、angels-burner-ore-crusher = ギア3 + 石の炉1 + 石5。

## 副産物が詰まるとき
- 炉は砕石（angels-stone-crushed）を受け付けない。クラッシャーの後に木箱を挟んで、粉砕鉱石と砕石をまとめて受け、炉へは粉砕鉱石だけを渡す。
- 砕石は 15枚/分 のとき約11個/分 溜まる。木箱（3,200個）が満杯になるまで約2.7時間あるので、それまでに回収して stone に変える。
- 詰まったら、次の順に確認する。
  1. クラッシャーの出力スロット
  2. 箱の空き
  3. 燃料切れ

## 未検証
- 建物の速度と燃料消費（list_prototypes では大きさしか取れない）
- 箱の手前の burner-inserter の能力（目標時に約34個/分を運ぶ必要があり、余裕がない見込み）

どちらも、鉄板ライン（issue #3）の実測結果で確認する。

## 次のステップ
`plans/copper-early.md` を作る（CLAUDE.md のテンプレに従う）。そのとき angels-ore3 鉱床の濃い部分を observe で探す。
