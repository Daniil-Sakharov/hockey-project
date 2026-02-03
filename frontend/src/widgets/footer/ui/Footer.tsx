export function Footer() {
  const year = new Date().getFullYear()

  return (
    <footer className="relative border-t border-white/5 bg-[#0a0e1a] py-12">
      <div className="mx-auto max-w-6xl px-4">
        <div className="flex flex-col items-center justify-between gap-6 md:flex-row">
          {/* Logo and description */}
          <div className="text-center md:text-left">
            <div className="mb-2 text-lg font-bold">
              <span className="text-gradient">Hockey</span>
              <span className="text-white">Stats</span>
            </div>
            <p className="text-sm text-gray-500">
              Аналитика российского хоккея
            </p>
          </div>

          {/* Links */}
          <div className="flex gap-8">
            <a
              href="#"
              className="text-sm text-gray-500 transition-colors hover:text-[#00d4ff]"
            >
              О проекте
            </a>
            <a
              href="#"
              className="text-sm text-gray-500 transition-colors hover:text-[#00d4ff]"
            >
              Контакты
            </a>
            <a
              href="#"
              className="text-sm text-gray-500 transition-colors hover:text-[#00d4ff]"
            >
              API
            </a>
          </div>
        </div>

        {/* Copyright */}
        <div className="mt-8 border-t border-white/5 pt-8 text-center">
          <p className="text-sm text-gray-600">
            © {year} HockeyStats. Все права защищены.
          </p>
        </div>
      </div>

      {/* Decorative gradient line */}
      <div className="absolute left-0 right-0 top-0 h-px bg-gradient-to-r from-transparent via-[#00d4ff]/30 to-transparent" />
    </footer>
  )
}
