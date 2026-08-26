#!/usr/bin/env bash
# 今日の作業ログを docs/logs/ にテンプレートから作成する。
#
# 使い方:
#   scripts/new-log.sh            # 今日の日付でログを作る
#   scripts/new-log.sh 2026-09-01 # 日付を指定して作る
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
log_dir="${repo_root}/docs/logs"
template="${log_dir}/TEMPLATE.md"

date_str="${1:-$(date +%F)}"
target="${log_dir}/${date_str}.md"

if [[ ! -f "${template}" ]]; then
  echo "テンプレートが見つかりません: ${template}" >&2
  exit 1
fi

if [[ -f "${target}" ]]; then
  echo "既に存在します: ${target}"
  exit 0
fi

sed "s/{{DATE}}/${date_str}/g" "${template}" > "${target}"
echo "作成しました: ${target}"
