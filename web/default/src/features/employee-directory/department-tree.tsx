/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/
import { ChevronRight, Folder, FolderOpen } from 'lucide-react'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'

import { Button } from '@/components/ui/button'
import { cn } from '@/lib/utils'

import type { DepartmentTreeNode } from './types'

type DepartmentTreeProps = {
  nodes: DepartmentTreeNode[]
  selectedId: number
  search: string
  onSelect: (department: DepartmentTreeNode) => void
}

function DepartmentNode(props: {
  node: DepartmentTreeNode
  selectedId: number
  forceOpen: boolean
  onSelect: (department: DepartmentTreeNode) => void
}) {
  const { t } = useTranslation()
  const [open, setOpen] = useState(props.node.level === 0)
  const children = props.node.children ?? []
  const expanded = props.forceOpen || open
  const hasChildren = children.length > 0
  const departmentName =
    props.node.id === -1 ? t('Unassigned') : props.node.name

  return (
    <li>
      <div
        className={cn(
          'group flex h-9 items-center gap-1 rounded-md px-1',
          props.selectedId === props.node.id &&
            'bg-accent text-accent-foreground'
        )}
        style={{ paddingLeft: `${props.node.level * 12 + 4}px` }}
      >
        <Button
          type='button'
          variant='ghost'
          size='icon'
          className='size-7 shrink-0'
          disabled={!hasChildren}
          aria-label={
            expanded ? t('Collapse department') : t('Expand department')
          }
          onClick={() => setOpen((value) => !value)}
        >
          <ChevronRight
            className={cn(
              'size-4 transition-transform',
              expanded && 'rotate-90',
              !hasChildren && 'opacity-0'
            )}
          />
        </Button>
        <button
          type='button'
          className='flex min-w-0 flex-1 items-center gap-2 text-left text-sm'
          onClick={() => props.onSelect(props.node)}
        >
          {expanded && hasChildren ? (
            <FolderOpen className='text-muted-foreground size-4 shrink-0' />
          ) : (
            <Folder className='text-muted-foreground size-4 shrink-0' />
          )}
          <span className='truncate'>{departmentName}</span>
          <span className='text-muted-foreground ml-auto shrink-0 text-xs tabular-nums'>
            {props.node.total_employees}
          </span>
        </button>
      </div>
      {hasChildren && expanded && (
        <ul>
          {children.map((child) => (
            <DepartmentNode
              key={child.id}
              node={child}
              selectedId={props.selectedId}
              forceOpen={props.forceOpen}
              onSelect={props.onSelect}
            />
          ))}
        </ul>
      )}
    </li>
  )
}

function filterDepartments(
  nodes: DepartmentTreeNode[],
  search: string,
  unassignedLabel: string
): DepartmentTreeNode[] {
  if (!search) return nodes
  return nodes.flatMap((node) => {
    const children = filterDepartments(
      node.children ?? [],
      search,
      unassignedLabel
    )
    const name = node.id === -1 ? unassignedLabel : node.name
    if (name.toLowerCase().includes(search) || children.length > 0) {
      return [{ ...node, children }]
    }
    return []
  })
}

export function DepartmentTree(props: DepartmentTreeProps) {
  const { t } = useTranslation()
  const search = props.search.trim().toLowerCase()
  const nodes = filterDepartments(props.nodes, search, t('Unassigned'))

  return (
    <ul className='space-y-0.5'>
      {nodes.map((node) => (
        <DepartmentNode
          key={node.id}
          node={node}
          selectedId={props.selectedId}
          forceOpen={search.length > 0}
          onSelect={props.onSelect}
        />
      ))}
    </ul>
  )
}
