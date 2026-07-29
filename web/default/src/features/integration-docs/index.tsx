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
import { Activity, BarChart3, Clock3, KeyRound } from 'lucide-react'
import { useTranslation } from 'react-i18next'

import { PublicLayout } from '@/components/layout'

type PeriodRow = {
  value: string
  description: string
  buckets: string
}

type FieldRow = {
  field: string
  description: string
}

const requestExample = `GET /api/usage/token/employee?employee_no=zt64003&period=7d HTTP/1.1
Authorization: Bearer sk-xxx`

const headerExample = `GET /api/usage/token/employee?period=realtime HTTP/1.1
Authorization: Bearer sk-xxx
X-Employee-No: zt64003`

const responseExample = `{
  "success": true,
  "message": "",
  "data": {
    "object": "employee_token_usage",
    "token_id": 12,
    "employee_no": "zt64003",
    "period": "7d",
    "currency": "USD",
    "start_timestamp": 1784563200,
    "end_timestamp": 1785123600,
    "summary": {
      "request_count": 18,
      "quota": 0.00252,
      "prompt_tokens": 3200,
      "completion_tokens": 1800,
      "total_tokens": 5000
    },
    "buckets": [
      {
        "label": "2026-07-15",
        "start_timestamp": 1784563200,
        "end_timestamp": 1784649600,
        "usage": {
          "request_count": 3,
          "quota": 0.00052,
          "prompt_tokens": 800,
          "completion_tokens": 420,
          "total_tokens": 1220
        }
      }
    ]
  }
}`

function CodeBlock(props: { children: string }) {
  return (
    <pre className='border-border/60 bg-muted/35 overflow-x-auto rounded-lg border px-4 py-3 text-xs leading-relaxed'>
      <code>{props.children}</code>
    </pre>
  )
}

export function IntegrationDocs() {
  const { t } = useTranslation()

  const periodRows: PeriodRow[] = [
    {
      value: 'realtime',
      description: t('Last 60 seconds for live dashboards.'),
      buckets: t('No daily buckets.'),
    },
    {
      value: 'today',
      description: t('From local midnight to the current time.'),
      buckets: t('No daily buckets.'),
    },
    {
      value: '7d',
      description: t('The last 7 calendar days including today.'),
      buckets: t('Returns daily buckets.'),
    },
    {
      value: '30d',
      description: t('The last 30 calendar days including today.'),
      buckets: t('Returns daily buckets.'),
    },
    {
      value: 'month',
      description: t('From the first day of the current month to now.'),
      buckets: t('Returns daily buckets.'),
    },
  ]

  const fieldRows: FieldRow[] = [
    {
      field: 'request_count',
      description: t('Number of successful consumption log records.'),
    },
    {
      field: 'quota',
      description: t('Amount consumed by the employee under this key, in the response currency.'),
    },
    {
      field: 'prompt_tokens',
      description: t('Input tokens recorded in the consumption log.'),
    },
    {
      field: 'completion_tokens',
      description: t('Output tokens recorded in the consumption log.'),
    },
    {
      field: 'total_tokens',
      description: t('prompt_tokens plus completion_tokens.'),
    },
  ]

  return (
    <PublicLayout>
      <main className='mx-auto w-full max-w-5xl px-4 py-12 md:py-16'>
        <section className='mb-10 space-y-5'>
          <div className='border-border/60 bg-muted/20 inline-flex items-center gap-2 rounded-full border px-3 py-1 text-xs text-muted-foreground'>
            <KeyRound className='size-3.5' />
            {t('Business System Integration')}
          </div>
          <div className='space-y-4'>
            <h1 className='text-3xl font-semibold tracking-tight md:text-4xl'>
              {t('Integration Docs')}
            </h1>
            <p className='text-muted-foreground max-w-3xl text-sm leading-7 md:text-base'>
              {t(
                'Business systems can use their assigned API key to query token usage for a specific employee number under that key.'
              )}
            </p>
          </div>
        </section>

        <section className='grid gap-4 md:grid-cols-3'>
          <div className='border-border/60 rounded-lg border p-4'>
            <KeyRound className='mb-3 size-5 text-blue-500' />
            <h2 className='mb-1 text-sm font-medium'>{t('Authentication')}</h2>
            <p className='text-muted-foreground text-sm leading-6'>
              {t('Use the business system API key in the Bearer token.')}
            </p>
          </div>
          <div className='border-border/60 rounded-lg border p-4'>
            <Activity className='mb-3 size-5 text-emerald-500' />
            <h2 className='mb-1 text-sm font-medium'>{t('Employee Scope')}</h2>
            <p className='text-muted-foreground text-sm leading-6'>
              {t('Pass employee_no as a query parameter or X-Employee-No header.')}
            </p>
          </div>
          <div className='border-border/60 rounded-lg border p-4'>
            <Clock3 className='mb-3 size-5 text-amber-500' />
            <h2 className='mb-1 text-sm font-medium'>{t('Time Ranges')}</h2>
            <p className='text-muted-foreground text-sm leading-6'>
              {t('Query realtime, daily, 7-day, 30-day, or monthly usage.')}
            </p>
          </div>
        </section>

        <section className='mt-10 space-y-4'>
          <h2 className='flex items-center gap-2 text-lg font-semibold'>
            <BarChart3 className='size-5' />
            {t('Employee Token Usage API')}
          </h2>
          <div className='border-border/60 overflow-hidden rounded-lg border'>
            <div className='grid gap-px bg-border/60 md:grid-cols-[180px_1fr]'>
              <div className='bg-background px-4 py-3 text-sm font-medium'>
                {t('Method')}
              </div>
              <div className='bg-background px-4 py-3 text-sm'>GET</div>
              <div className='bg-background px-4 py-3 text-sm font-medium'>
                {t('Path')}
              </div>
              <div className='bg-background px-4 py-3 font-mono text-sm'>
                /api/usage/token/employee
              </div>
              <div className='bg-background px-4 py-3 text-sm font-medium'>
                {t('Authentication')}
              </div>
              <div className='bg-background px-4 py-3 font-mono text-sm'>
                Authorization: Bearer sk-xxx
              </div>
            </div>
          </div>
        </section>

        <section className='mt-10 space-y-4'>
          <h2 className='text-lg font-semibold'>{t('Request Examples')}</h2>
          <div className='grid gap-4 lg:grid-cols-2'>
            <div className='space-y-2'>
              <p className='text-muted-foreground text-sm'>
                {t('Pass employee_no in query parameters.')}
              </p>
              <CodeBlock>{requestExample}</CodeBlock>
            </div>
            <div className='space-y-2'>
              <p className='text-muted-foreground text-sm'>
                {t('Pass the employee number in the request header.')}
              </p>
              <CodeBlock>{headerExample}</CodeBlock>
            </div>
          </div>
        </section>

        <section className='mt-10 space-y-4'>
          <h2 className='text-lg font-semibold'>{t('Supported Periods')}</h2>
          <div className='border-border/60 overflow-hidden rounded-lg border'>
            <table className='w-full text-left text-sm'>
              <thead className='bg-muted/40 text-muted-foreground'>
                <tr>
                  <th className='px-4 py-3 font-medium'>{t('Value')}</th>
                  <th className='px-4 py-3 font-medium'>{t('Meaning')}</th>
                  <th className='px-4 py-3 font-medium'>{t('Buckets')}</th>
                </tr>
              </thead>
              <tbody className='divide-border divide-y'>
                {periodRows.map((row) => (
                  <tr key={row.value}>
                    <td className='px-4 py-3 font-mono'>{row.value}</td>
                    <td className='px-4 py-3 text-muted-foreground'>
                      {row.description}
                    </td>
                    <td className='px-4 py-3 text-muted-foreground'>
                      {row.buckets}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>

        <section className='mt-10 space-y-4'>
          <h2 className='text-lg font-semibold'>{t('Response Fields')}</h2>
          <div className='border-border/60 overflow-hidden rounded-lg border'>
            <table className='w-full text-left text-sm'>
              <thead className='bg-muted/40 text-muted-foreground'>
                <tr>
                  <th className='px-4 py-3 font-medium'>{t('Field')}</th>
                  <th className='px-4 py-3 font-medium'>{t('Description')}</th>
                </tr>
              </thead>
              <tbody className='divide-border divide-y'>
                {fieldRows.map((row) => (
                  <tr key={row.field}>
                    <td className='px-4 py-3 font-mono'>{row.field}</td>
                    <td className='px-4 py-3 text-muted-foreground'>
                      {row.description}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </section>

        <section className='mt-10 space-y-4'>
          <h2 className='text-lg font-semibold'>{t('Response Example')}</h2>
          <CodeBlock>{responseExample}</CodeBlock>
        </section>

        <section className='mt-10 rounded-lg border border-border/60 bg-muted/20 p-4'>
          <h2 className='mb-2 text-sm font-medium'>{t('Integration Notes')}</h2>
          <ul className='text-muted-foreground list-disc space-y-2 pl-5 text-sm leading-6'>
            <li>{t('The API key only sees usage generated under that key.')}</li>
            <li>
              {t(
                'employee_no is treated as an external identifier and does not require employee import.'
              )}
            </li>
            <li>
              {t(
                'Only successful consumption logs are included; error logs are excluded.'
              )}
            </li>
          </ul>
        </section>
      </main>
    </PublicLayout>
  )
}
