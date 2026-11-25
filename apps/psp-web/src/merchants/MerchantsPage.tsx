import { useEffect, useState, type FormEvent } from 'react'
import { apiFetch } from '../api/apiFetch'

type Merchant = {
  id: number
  name: string
  email: string
  contactFirstName: string
  contactLastName: string
  createdAt: string
}

type CreateMerchantResponse = {
  merchant: Merchant
  temporaryPassword: string
}

const emptyForm = { name: '', email: '', contactFirstName: '', contactLastName: '' }

export function MerchantsPage() {
  const [merchants, setMerchants] = useState<Merchant[]>([])
  const [form, setForm] = useState(emptyForm)
  const [createdMerchant, setCreatedMerchant] = useState<CreateMerchantResponse | null>(null)
  const [errorMessage, setErrorMessage] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  useEffect(() => {
    apiFetch<Merchant[]>('/merchants')
      .then(setMerchants)
      .catch((error: Error) => setErrorMessage(error.message))
  }, [])

  async function handleSubmit(event: FormEvent) {
    event.preventDefault()
    setErrorMessage('')
    setCreatedMerchant(null)
    setIsSubmitting(true)

    try {
      const response = await apiFetch<CreateMerchantResponse>('/merchants', { method: 'POST', body: form })
      setCreatedMerchant(response)
      setMerchants([response.merchant, ...merchants])
      setForm(emptyForm)
    } catch (error) {
      setErrorMessage((error as Error).message)
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <section className="w-full max-w-3xl space-y-6">
      <form onSubmit={handleSubmit} className="space-y-3 rounded border border-slate-200 bg-white p-4">
        <h2 className="text-lg font-semibold">Register merchant</h2>
        <div className="grid grid-cols-2 gap-3">
          <TextField label="Merchant name" value={form.name} onChange={(name) => setForm({ ...form, name })} />
          <TextField label="Email (login)" type="email" value={form.email} onChange={(email) => setForm({ ...form, email })} />
          <TextField
            label="Contact first name"
            value={form.contactFirstName}
            onChange={(contactFirstName) => setForm({ ...form, contactFirstName })}
          />
          <TextField
            label="Contact last name"
            value={form.contactLastName}
            onChange={(contactLastName) => setForm({ ...form, contactLastName })}
          />
        </div>
        <button
          type="submit"
          disabled={isSubmitting}
          className="rounded bg-slate-800 px-4 py-2 text-sm text-white hover:bg-slate-700 disabled:opacity-50"
        >
          Register
        </button>
        {errorMessage && <p className="text-sm text-red-600">{errorMessage}</p>}
      </form>

      {createdMerchant && (
        <div className="rounded border border-amber-300 bg-amber-50 p-4 text-sm">
          <p>
            Merchant <strong>{createdMerchant.merchant.name}</strong> registered. Temporary password:
          </p>
          <p className="my-2 font-mono text-base">{createdMerchant.temporaryPassword}</p>
          <p>
            Give it to the merchant securely. It is shown only once. On first login the merchant must change it and
            set up an authenticator app.
          </p>
        </div>
      )}

      <table className="w-full rounded border border-slate-200 bg-white text-left text-sm">
        <thead className="bg-slate-100">
          <tr>
            <th className="p-2">Name</th>
            <th className="p-2">Email</th>
            <th className="p-2">Contact</th>
            <th className="p-2">Registered</th>
          </tr>
        </thead>
        <tbody>
          {merchants.map((merchant) => (
            <tr key={merchant.id} className="border-t border-slate-200">
              <td className="p-2">{merchant.name}</td>
              <td className="p-2">{merchant.email}</td>
              <td className="p-2">
                {merchant.contactFirstName} {merchant.contactLastName}
              </td>
              <td className="p-2">{new Date(merchant.createdAt).toLocaleString()}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </section>
  )
}

type TextFieldProps = {
  label: string
  value: string
  onChange: (value: string) => void
  type?: string
}

function TextField({ label, value, onChange, type = 'text' }: TextFieldProps) {
  return (
    <label className="flex flex-col gap-1 text-sm">
      {label}
      <input
        required
        type={type}
        value={value}
        onChange={(event) => onChange(event.target.value)}
        className="rounded border border-slate-300 px-2 py-1"
      />
    </label>
  )
}
