import { useEffect, useState } from 'react'
import { paymentService } from '@/services/api'
import { toast } from 'react-toastify'
import { FileDown, Printer } from 'lucide-react'
import { PageHeader, Toolbar, SearchInput, Loading, EmptyState } from '@/components/common'

export default function PaymentsPage() {
  const [payments, setPayments] = useState<any[]>([])
  const [filtered, setFiltered] = useState<any[]>([])
  const [search, setSearch] = useState('')
  const [loading, setLoading] = useState(true)

  const load = () => {
    setLoading(true)
    paymentService.list()
      .then(r => { setPayments(r.data.data || []) })
      .catch(() => toast.error('Failed to load payments'))
      .finally(() => setLoading(false))
  }

  useEffect(() => { load() }, [])

  useEffect(() => {
    if (!search) { setFiltered(payments); return }
    const q = search.toLowerCase()
    setFiltered(payments.filter(p =>
      p.invoice?.invoice_number?.toLowerCase().includes(q) ||
      p.payment_method?.toLowerCase().includes(q)
    ))
  }, [search, payments])

  const fmt = (n: number, cur = 'IDR') => `${cur} ${Number(n).toLocaleString()}`
  const total = payments.reduce((sum, p) => sum + (p.amount || 0), 0)

  return (
    <div className="p-5">
      <PageHeader title="Payments" />

      <div className="grid grid-cols-3 gap-4 mb-4">
        <div className="bg-white rounded-lg border border-gray-200 p-4">
          <p className="text-xs text-gray-400 mb-1">Total Payments</p>
          <p className="text-xl font-semibold text-gray-900">{payments.length}</p>
        </div>
        <div className="bg-white rounded-lg border border-gray-200 p-4">
          <p className="text-xs text-gray-400 mb-1">Total Received</p>
          <p className="text-xl font-semibold text-green-600">{fmt(total)}</p>
        </div>
      </div>

      <Toolbar
        left={<span className="text-xs text-gray-400">{filtered.length} records</span>}
        right={
          <>
            <button className="btn btn-secondary"><FileDown size={12} />Excel</button>
            <button className="btn btn-secondary"><Printer size={12} />Print</button>
            <SearchInput value={search} onChange={setSearch} placeholder="Search invoice..." />
          </>
        }
      />

      <div className="table-container">
        {loading ? <Loading /> : (
          <table className="table">
            <thead>
              <tr><th>Invoice #</th><th>Amount</th><th>Method</th><th>Date</th><th>Note</th></tr>
            </thead>
            <tbody>
              {filtered.length === 0
                ? <tr><td colSpan={5}><EmptyState /></td></tr>
                : filtered.map(p => (
                  <tr key={p.id}>
                    <td className="font-medium text-blue-600">{p.invoice?.invoice_number || '-'}</td>
                    <td className="whitespace-nowrap font-medium text-green-600">{fmt(p.amount, p.currency)}</td>
                    <td className="text-gray-500 capitalize">{p.payment_method || '-'}</td>
                    <td className="text-gray-400 whitespace-nowrap">{p.payment_date ? new Date(p.payment_date).toLocaleDateString('id') : '-'}</td>
                    <td className="text-gray-400">{p.note || '-'}</td>
                  </tr>
                ))
              }
            </tbody>
          </table>
        )}
      </div>
    </div>
  )
}
