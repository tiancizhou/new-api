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
export type DirectoryStatus = 1 | 2

export type DepartmentTreeNode = {
  id: number
  name: string
  parent_id: number
  level: number
  status: DirectoryStatus
  external_dept_id: string
  employee_count: number
  total_employees: number
  children?: DepartmentTreeNode[]
}

export type Employee = {
  id: number
  employee_no: string
  external_employee_id: string
  name: string
  email?: string
  phone?: string
  department_id: number
  external_dept_id?: string
  department_name: string
  post_id?: string
  sex?: number
  status: DirectoryStatus
  created_at: number
  updated_at: number
}

export type EmployeePage = {
  page: number
  page_size: number
  total: number
  items: Employee[]
}

export type ApiResponse<T> = {
  success: boolean
  message?: string
  data?: T
}

export type EmployeeQuery = {
  page: number
  pageSize: number
  keyword: string
  departmentId: number
  includeSubdepartments: boolean
  status: 'all' | '1' | '2'
}
