type SelectorProps<T> = {
  label?: string
  value: T | null
  setValue: (value: T | null) => void
  options: { label: string; value: T }[]
  allLabel?: string
  className?: string
}

export const Selector = <T,>({
  label = "Select",
  value,
  setValue,
  options,
  allLabel = "All",
  className = "col-3",
}: SelectorProps<T>) => {
  return (
    <div className={className}>
      <label className="form-label">{label}</label>
      <select
        className="form-select"
        value={value as any}
        onChange={(e) => {
          const val = e.target.value
          // Если значение пустое — "All"
          if (val === "") {
            setValue(null)
            return
          }

          // Пробуем привести строку обратно к типу T
          const found = options?.find((o) => String(o.value) === val)
          if (found) {
            setValue(found.value)
          }
        }}
      >
        <option value="">{allLabel}</option>
        {options?.map((opt) => (
          <option key={String(opt.value)} value={String(opt.value)}>
            {opt.label}
          </option>
        ))}
      </select>
    </div>
  )
}
