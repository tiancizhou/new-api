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
import { keepPreviousData, useQuery } from '@tanstack/react-query'
import {
  Building2,
  RefreshCw,
  Search,
  UserRound,
  UsersRound,
} from 'lucide-react'
import { useDeferredValue, useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'

import { SectionPageLayout } from '@/components/layout'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  Empty,
  EmptyDescription,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
} from '@/components/ui/empty'
import { Input } from '@/components/ui/input'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Skeleton } from '@/components/ui/skeleton'
import { Switch } from '@/components/ui/switch'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'

import { getDepartmentTree, getEmployees } from './api'
import { DepartmentTree } from './department-tree'
import type { DepartmentTreeNode, EmployeeQuery } from './types'

const PAGE_SIZE = 20

export function EmployeeDirectory() {
  const { t } = useTranslation()
  const [departmentId, setDepartmentId] = useState(0)
  const [departmentName, setDepartmentName] = useState('')
  const [departmentSearch, setDepartmentSearch] = useState('')
  const [employeeSearch, setEmployeeSearch] = useState('')
  const [includeSubdepartments, setIncludeSubdepartments] = useState(true)
  const [status, setStatus] = useState<EmployeeQuery['status']>('1')
  const [page, setPage] = useState(1)
  const deferredEmployeeSearch = useDeferredValue(employeeSearch.trim())

  const departmentsQuery = useQuery({
    queryKey: ['employee-directory', 'departments'],
    queryFn: getDepartmentTree,
  })

  const employeeQuery = useMemo<EmployeeQuery>(
    () => ({
      page,
      pageSize: PAGE_SIZE,
      keyword: deferredEmployeeSearch,
      departmentId,
      includeSubdepartments,
      status,
    }),
    [page, deferredEmployeeSearch, departmentId, includeSubdepartments, status]
  )

  const employeesQuery = useQuery({
    queryKey: ['employee-directory', 'employees', employeeQuery],
    queryFn: () => getEmployees(employeeQuery),
    placeholderData: keepPreviousData,
  })

  const departments = departmentsQuery.data ?? []
  const totalEmployees = departments.reduce(
    (total, department) => total + department.total_employees,
    0
  )
  const employeePage = employeesQuery.data
  const totalPages = Math.max(
    1,
    Math.ceil((employeePage?.total ?? 0) / PAGE_SIZE)
  )

  const selectDepartment = (department: DepartmentTreeNode) => {
    setDepartmentId(department.id)
    setDepartmentName(department.id === -1 ? t('Unassigned') : department.name)
    setPage(1)
  }

  return (
    <SectionPageLayout fixedContent>
      <SectionPageLayout.Title>
        {t('Employee Management')}
      </SectionPageLayout.Title>
      <SectionPageLayout.Actions>
        <Button
          variant='outline'
          size='sm'
          onClick={() => {
            departmentsQuery.refetch()
            employeesQuery.refetch()
          }}
          disabled={departmentsQuery.isFetching || employeesQuery.isFetching}
        >
          <RefreshCw className='size-4' />
          {t('Refresh')}
        </Button>
      </SectionPageLayout.Actions>
      <SectionPageLayout.Content>
        <div className='grid h-full min-h-0 grid-cols-1 gap-3 lg:grid-cols-[280px_minmax(0,1fr)]'>
          <aside className='bg-background flex min-h-0 flex-col overflow-hidden rounded-md border lg:h-full'>
            <div className='flex h-12 shrink-0 items-center justify-between border-b px-3'>
              <div className='flex items-center gap-2 text-sm font-semibold'>
                <Building2 className='size-4' />
                {t('Departments')}
              </div>
              <Badge variant='secondary' className='tabular-nums'>
                {totalEmployees}
              </Badge>
            </div>
            <div className='border-b p-2'>
              <div className='relative'>
                <Search className='text-muted-foreground pointer-events-none absolute top-1/2 left-2.5 size-4 -translate-y-1/2' />
                <Input
                  value={departmentSearch}
                  onChange={(event) => setDepartmentSearch(event.target.value)}
                  placeholder={t('Search departments')}
                  className='h-9 pl-8'
                />
              </div>
            </div>
            <button
              type='button'
              className={`mx-2 mt-2 flex h-9 items-center gap-2 rounded-md px-3 text-left text-sm ${departmentId === 0 ? 'bg-accent text-accent-foreground' : 'hover:bg-muted'}`}
              onClick={() => {
                setDepartmentId(0)
                setDepartmentName('')
                setPage(1)
              }}
            >
              <UsersRound className='text-muted-foreground size-4' />
              <span className='flex-1'>{t('All departments')}</span>
              <span className='text-muted-foreground text-xs tabular-nums'>
                {totalEmployees}
              </span>
            </button>
            <ScrollArea className='h-72 min-h-0 p-2 lg:h-auto lg:flex-1'>
              {departmentsQuery.isLoading && (
                <div className='space-y-2 p-1'>
                  {Array.from({ length: 8 }, (_, index) => (
                    <Skeleton key={index} className='h-8 w-full' />
                  ))}
                </div>
              )}
              {!departmentsQuery.isLoading && departments.length > 0 && (
                <DepartmentTree
                  nodes={departments}
                  selectedId={departmentId}
                  search={departmentSearch}
                  onSelect={selectDepartment}
                />
              )}
              {!departmentsQuery.isLoading && departments.length === 0 && (
                <p className='text-muted-foreground p-4 text-center text-sm'>
                  {t('No departments found')}
                </p>
              )}
            </ScrollArea>
          </aside>

          <section className='bg-background flex min-h-0 flex-col overflow-hidden rounded-md border'>
            <div className='flex shrink-0 flex-col gap-3 border-b p-3 xl:flex-row xl:items-center xl:justify-between'>
              <div className='min-w-0'>
                <h3 className='truncate text-sm font-semibold'>
                  {departmentName || t('All employees')}
                </h3>
                <p className='text-muted-foreground mt-0.5 text-xs'>
                  {t('{{count}} employees', {
                    count: employeePage?.total ?? 0,
                  })}
                </p>
              </div>
              <div className='flex flex-wrap items-center gap-2'>
                <div className='relative min-w-52 flex-1 sm:flex-none'>
                  <Search className='text-muted-foreground pointer-events-none absolute top-1/2 left-2.5 size-4 -translate-y-1/2' />
                  <Input
                    value={employeeSearch}
                    onChange={(event) => {
                      setEmployeeSearch(event.target.value)
                      setPage(1)
                    }}
                    placeholder={t('Search by employee number or name')}
                    className='h-9 w-full pl-8 sm:w-64'
                  />
                </div>
                <select
                  value={status}
                  onChange={(event) => {
                    setStatus(event.target.value as EmployeeQuery['status'])
                    setPage(1)
                  }}
                  className='border-input bg-background h-9 rounded-md border px-3 text-sm'
                  aria-label={t('Status')}
                >
                  <option value='1'>{t('Enabled')}</option>
                  <option value='2'>{t('Disabled')}</option>
                  <option value='all'>{t('All statuses')}</option>
                </select>
                {departmentId > 0 && (
                  <label className='flex h-9 items-center gap-2 rounded-md border px-3 text-sm whitespace-nowrap'>
                    <Switch
                      checked={includeSubdepartments}
                      onCheckedChange={(checked) => {
                        setIncludeSubdepartments(checked)
                        setPage(1)
                      }}
                    />
                    {t('Include subdepartments')}
                  </label>
                )}
              </div>
            </div>

            <div className='min-h-0 flex-1 overflow-auto'>
              {employeesQuery.isLoading && (
                <div className='space-y-3 p-4'>
                  {Array.from({ length: 8 }, (_, index) => (
                    <Skeleton key={index} className='h-12 w-full' />
                  ))}
                </div>
              )}
              {!employeesQuery.isLoading &&
                employeePage &&
                employeePage.items.length > 0 && (
                  <Table>
                    <TableHeader className='bg-muted/60 sticky top-0 z-10'>
                      <TableRow>
                        <TableHead className='w-36'>
                          {t('Employee Number')}
                        </TableHead>
                        <TableHead>{t('Name')}</TableHead>
                        <TableHead>{t('Department')}</TableHead>
                        <TableHead>{t('Email')}</TableHead>
                        <TableHead>{t('Phone')}</TableHead>
                        <TableHead className='w-24'>{t('Status')}</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {employeePage.items.map((employee) => (
                        <TableRow key={employee.id}>
                          <TableCell>
                            <div className='flex items-center gap-2 font-mono font-medium'>
                              <UserRound className='text-muted-foreground size-4' />
                              {employee.employee_no}
                            </div>
                          </TableCell>
                          <TableCell className='font-medium'>
                            {employee.name}
                          </TableCell>
                          <TableCell>
                            {employee.department_name || '-'}
                          </TableCell>
                          <TableCell className='max-w-64 truncate'>
                            {employee.email || '-'}
                          </TableCell>
                          <TableCell>{employee.phone || '-'}</TableCell>
                          <TableCell>
                            <Badge
                              variant={
                                employee.status === 1 ? 'default' : 'secondary'
                              }
                            >
                              {employee.status === 1
                                ? t('Enabled')
                                : t('Disabled')}
                            </Badge>
                          </TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                )}
              {!employeesQuery.isLoading && !employeePage?.items.length && (
                <Empty className='h-full border-0'>
                  <EmptyHeader>
                    <EmptyMedia variant='icon'>
                      <UserRound />
                    </EmptyMedia>
                    <EmptyTitle>{t('No employees found')}</EmptyTitle>
                    <EmptyDescription>
                      {t('Try adjusting the department or search filters.')}
                    </EmptyDescription>
                  </EmptyHeader>
                </Empty>
              )}
            </div>

            <div className='flex h-12 shrink-0 items-center justify-between border-t px-3'>
              <span className='text-muted-foreground text-xs tabular-nums'>
                {t('Page {{page}} of {{total}}', { page, total: totalPages })}
              </span>
              <div className='flex items-center gap-1'>
                <Button
                  variant='outline'
                  size='sm'
                  disabled={page <= 1}
                  onClick={() => setPage((value) => value - 1)}
                >
                  {t('Previous')}
                </Button>
                <Button
                  variant='outline'
                  size='sm'
                  disabled={page >= totalPages}
                  onClick={() => setPage((value) => value + 1)}
                >
                  {t('Next')}
                </Button>
              </div>
            </div>
          </section>
        </div>
      </SectionPageLayout.Content>
    </SectionPageLayout>
  )
}
