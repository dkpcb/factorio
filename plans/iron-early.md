# iron-early: 序盤の鉄板ライン

## 目的
サフィライト（angels-ore1）から鉄板（iron-plate）を作る。全部バーナー機で、電力網なしで動く序盤用ライン。

## surface
nauvis

## 範囲(座標)
- 外枠: left_top (56, 2) 〜 right_bottom (60, 14)
- 鉱床: angels-ore1 パッチ西寄りの最高濃度部（tick 26116 時点の observe 値）
  - パッチ範囲: 概ね x52〜90, y-15〜26（半径160内で829タイル）
  - 最高濃度帯: x56〜59, y2〜5（1タイル約1,700〜1,980）
- 近くの資源・障害物
  - coal: 外枠の北側 x56〜61, y-3〜0 で接している（パッチは北へ y-24 付近まで広がる）
  - 外枠を含む (54,0)〜(62,16) に木・岩・崖・自軍エンティティ・ゴーストはなし
  - 最寄りの木は tree-08-red (75.6, 11.0)。外枠の外なので撤去不要

## 目標生産量
- iron-plate 15/分
- 設計上の能力 17.1/分（ボトルネックは炉1基）

## 使用レシピ（list_recipes値）
| レシピ | カテゴリ | energy | 入力 | 出力 |
|---|---|---|---|---|
| angels-ore1-crushed | angels-ore-refining-t1 | 1 | angels-ore1 ×2 | angels-ore1-crushed ×2, angels-stone-crushed ×1 |
| iron-plate | smelting | 7 | angels-ore1-crushed ×3 | iron-plate ×2 |
| （副産物の処理用）angels-stone-from-crushed-stone | crafting | 0.5 | angels-stone-crushed ×2 | stone ×1 |

注: 炉は angels-stone-crushed を受け付けない（crafting カテゴリのため）。

建物のレシピ（すべて enabled=true）

| レシピ | 入力 |
|---|---|
| burner-mining-drill | iron-gear-wheel ×3, stone-furnace ×1, iron-plate ×3 |
| angels-burner-ore-crusher | iron-gear-wheel ×3, stone-furnace ×1, stone ×5 |
| stone-furnace | stone ×5 |
| burner-inserter | iron-gear-wheel ×1, iron-plate ×1 |
| wooden-chest | wood ×2 |
| iron-gear-wheel | iron-plate ×2 |

## 必要建物数・資材

### 計算（15/分）
| 工程 | 必要量 | 1台の能力 | 台数 | 稼働率 |
|---|---|---|---|---|
| 精錬 stone-furnace (speed 1) | crushed 22.5/分 → plate 15/分 | 2/7s = 17.1 plate/分 | 1 | 88% |
| 粉砕 angels-burner-ore-crusher (speed 1) | ore 22.5/分 | 2 ore/1s = 120 ore/分 | 1 | 19% |
| 採掘 burner-mining-drill (mining_speed 0.25, mining_time 1) | ore 22.5/分 | 15 ore/分 | 2 | 75% |

- 速度・消費電力・インサーター速度は list_prototypes では取れない（返るのは collision_box だけ）。この表の値は旧計画のもので、未検証。建設後に inspect_entity で確認する。
- サイズは確認済み: burner-mining-drill 2×2、angels-burner-ore-crusher 3×3
- 副産物 angels-stone-crushed: 目標量で 11.25/分、ドリル全力時は最大 15/分
- 使えるのはバーナー機のみ（electric-mining-drill / angels-ore-crusher / steel-furnace などはまだ未解禁）

### 配置（北が -y、下流は南へ）
| # | エンティティ | 中心座標 | 占有タイル | 向き・役割 | レシピ |
|---|---|---|---|---|---|
| D1 | burner-mining-drill | (57, 3) | x56-57, y2-3 | 南向き。(57,4) に排出し、クラッシャーへ入る | - |
| D2 | burner-mining-drill | (59, 3) | x58-59, y2-3 | 南向き。(59,4) に排出し、クラッシャーへ入る | - |
| C  | angels-burner-ore-crusher | (58.5, 5.5) | x57-59, y4-6 | - | angels-ore1-crushed |
| I1 | burner-inserter | (58.5, 7.5) | (58,7) | クラッシャー(北) → 箱B1(南) | - |
| B1 | wooden-chest | (58.5, 8.5) | (58,8) | 粉砕鉱石と砕石の混合バッファ | - |
| I2 | burner-inserter | (58.5, 9.5) | (58,9) | 箱B1(北) → 炉(南)。炉が受け付ける粉砕鉱石だけを運ぶ | - |
| F  | stone-furnace | (59, 11) | x58-59, y10-11 | - | iron-plate（自動） |
| I3 | burner-inserter | (58.5, 12.5) | (58,12) | 炉(北) → 箱B2(南) | - |
| B2 | wooden-chest | (58.5, 13.5) | (58,13) | 鉄板の出力 | - |

- ドリルの向きは南（2.x の direction=8）。排出位置は北向きの (-0.3,-1.3) を180°回した (+0.3,+1.3) と想定している。preview_plan で、D1/D2 の排出先がクラッシャー内にあることを確認する。
- インサーターは「北から拾って南へ置く」。direction の値は create_plan / preview_plan で確認する。
- ドリル下の鉱石量: D1 = 1680+1861+1835+1965 = 7,341、D2 = 1904+1797+1980+1853 = 7,534。どちらも 15/分 で約8.2時間もつ。
- クラッシャーと B1〜B2 は鉱石の上に乗る（約1,400〜1,900/タイル）。後で撤去すれば採掘できる。
- 砕石はクラッシャーから直接捨てずに B1 で受ける。こうすると粉砕鉱石を失わない。B1 には砕石と余った粉砕鉱石（約4/分）が溜まる。
- B1 の容量は 16スロット×200 = 3200個。溜まる速さは最大約19/分なので、約2.7時間で満杯になる。それまでに砕石を回収する（手作業で stone に変換できる）。
- リスク: I1 は目標時に約34個/分（crushed 22.5 + 砕石 11.25）を運ぶ必要がある。burner-inserter の箱間能力は約36/分と見込んでいて、余裕が少ない。I1 が混合比 2:1 で運ぶなら crushed は約24/分届き、炉の必要量 22.5/分は満たせる計算。足りなければ、I1 の隣 (57,7)→(57,8) に2本目のインサーターと箱を追加する。

### 資材
| 品目 | 数 |
|---|---|
| burner-mining-drill | 2 |
| angels-burner-ore-crusher | 1 |
| stone-furnace | 1（設置用） |
| burner-inserter | 3 |
| wooden-chest | 2 |
| coal（初期燃料） | 約200（ドリル各50、炉50、クラッシャー20、インサーター各5〜10） |

素材換算（手作りする場合）: iron-plate 33（ギア12個分の24 + 直接9）、stone 25（stone-furnace 4基分の20 + クラッシャーの5）、wood 4

### 燃料
- 消費電力（旧計画の値、未検証）: ドリル 150kW×2、クラッシャー 100kW（稼働率は約25%）、炉 90kW、インサーター少量。合計は約0.45〜0.5MW。
- coal の燃料値は 4MJ なので、消費は約7個/分。
- ドリル1台は coal 200個（1スタック）で約89分動く。
- coal パッチが外枠の北に接している（(57〜61, -3〜0)、ドリルから約3〜5タイル）。当面は手で補給する。coal 用のバーナードリルで自給する案は別の計画にする。

## 撤去対象
なし（外枠内に障害物はなく、自軍エンティティも 0件）

## 検証方法
1. apply_plan(mode="ghosts") の後、建設が終わったら get_action_status で完了を確認する。
2. inspect_entity で次を確認する。
   - D1/D2: 採掘中であること、燃料があること、出力が詰まっていないこと。mining_speed と消費電力の実値
   - C: レシピが angels-ore1-crushed であること、出力が詰まっていないこと。crafting_speed の実値
   - F: レシピが iron-plate で稼働していること
   - B1: 中身の増え方（砕石が溜まり続けていないか）
   - I1: 能力が足りているか（C の出力スロットに溜まり続けていないか）
3. B2 の iron-plate の数を2時点で測って /分 を出す。目標は 15/分以上（理論値 17.1/分）。
4. 届かない場合は、燃料切れ・I1 の能力・B1 の満杯・インサーターの向きを順に確認する。
