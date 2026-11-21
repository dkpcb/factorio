# red-science-early: 赤サイエンス 最小1列

## 目的
赤サイエンスパック（automation-science-pack）を組立機で作り、手で運ばずに研究所へ届ける。Issue #8 の完了条件「手で運ばずに研究所に届く」「研究が実際に進む」を、まず1列で確かめる。

- 歯車と基本回路基板は、最初は手で箱に補充する。自動化は稼働を確かめてから別に計画する。
- 建設はプレイヤーが手で行う。Claude はゴーストを置かない。

## surface
nauvis

## 範囲(座標)
- 外枠: left_top (35, -35) 〜 right_bottom (39, -26)
- 電源: 既存の小型電柱 (35.5, -37.5)。蒸気発電 (29.5〜33.5, -62.5〜-60.5) から南へ延びた電柱線の南端
- 障害物（tick 603699 の observe 値）: (32,-37)〜(41,-24) にあるのは coal 1タイル (40.5, -29.5, 残量25) だけ。外枠の外なので影響なし
- 鉄板ライン (77, -45) まで約40タイル、石炭の箱 (56.5, -32.5) まで約20タイル、銅板ライン (-37, 28) まで約90タイル

## 目標生産量
- 赤SP 約6個/分（組立機1台の上限。組立機1の速度 0.5 は旧マップ調査時の base 定義の値で、このゲームでは未検証）
- 最初の確認は「研究が進むこと」。数値は建設後に測る

## 使用レシピ（list_recipes値、tick 592812〜593700 時点、すべて enabled=true）
| レシピ | カテゴリ | 時間 | 入力 → 出力 |
|---|---|---|---|
| automation-science-pack | crafting | 5秒 | iron-gear-wheel ×1 + bob-basic-circuit-board ×1 → 1 |
| bob-basic-circuit-board | crafting / electronics | 1秒 | bob-wooden-board ×1 + copper-cable ×3 → 1 |
| lab | crafting | 2秒 | iron-gear-wheel ×10 + bob-basic-circuit-board ×10 + bob-basic-transport-belt ×4 → 1 |
| assembling-machine-1 | crafting | 0.5秒 | iron-gear-wheel ×6 + bob-basic-circuit-board ×2 + iron-plate ×24 → 1 |
| inserter | crafting | 0.5秒 | iron-gear-wheel ×1 + bob-basic-circuit-board ×1 + iron-plate ×1 → 1 |

## 配置（北が -y、流れは北 → 南）
| # | エンティティ | 中心座標 | 占有タイル | 役割 |
|---|---|---|---|---|
| B1 | wooden-chest | (37.5, -34.5) | (37, -35) | 入力箱。歯車と基本回路基板を手で補充する |
| I1 | inserter | (37.5, -33.5) | (37, -34) | B1(北) から取って組立機(南) へ置く |
| A1 | assembling-machine-1 | (37.5, -31.5) | x36-38, y-33〜-31 | レシピ automation-science-pack |
| I2 | inserter | (37.5, -29.5) | (37, -30) | 組立機(北) から取って研究所(南) へ置く |
| L1 | lab | (37.5, -27.5) | x36-38, y-29〜-27 | 研究所 |
| P1 | small-electric-pole | (35.5, -31.5) | (35, -32) | 既存の電柱 (35.5,-37.5) から6タイル。I1・A1・I2 に給電 |
| P2 | small-electric-pole | (35.5, -27.5) | (35, -28) | P1 から4タイル。L1 に給電 |

- 小型電柱の給電範囲は 5×5（base の値。このゲームでは未検証）。置くときに、各設備の電源アイコンで給電範囲に入っているかを確認する。
- インサーターは、黄色の矢印が南を向くように置く（北から取り、南へ置く）。
- 後で歯車・基板を自動化するときは、B1 を外して西側 (x≤34) に組立機を並べる。そのために西側を空けておく。

## 必要建物数・資材
| 品目 | 数 | 鉄板 | 銅板 | 木材 |
|---|---|---|---|---|
| lab | 1 | 26 | 15 | 10 |
| assembling-machine-1 | 1 | 36 | 3 | 2 |
| inserter | 2 | 6 | 3 | 2 |
| small-electric-pole | 2（1回の製作で2本） | 0 | 1 | 1 |
| wooden-chest | 1 | 0 | 0 | 2 |
| **合計** | | **68** | **22** | **17** |

- 試運転用に、赤SP 10個分の歯車10個（鉄板20）と基本回路基板10個（銅板15、木材10）も用意する。

## 撤去対象
なし

## 検証方法
1. 組み上がったら、研究を1つ始める（プレイヤーが選ぶ）
2. Claude が `observe`（範囲は外枠のみ、detailed=true）で、A1 のレシピと状態、L1 に赤SPが入っているかを確かめる
3. `research` の progress が 0 から増えていることを確かめる
4. 数分あけて、A1 の製作数（または B1 の減り方）から、赤SPが毎分何個出ているかを測る
