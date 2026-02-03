import { memo } from 'react'
import { motion } from 'framer-motion'
import { Check, Sparkles, Crown, Zap } from 'lucide-react'
import { cn } from '@/shared/lib/utils'

interface PricingFeature {
  text: string
  included: boolean
}

interface PricingPlan {
  name: string
  price: string
  period: string
  description: string
  features: PricingFeature[]
  buttonText: string
  popular?: boolean
  icon: typeof Zap
  color: 'cyan' | 'purple' | 'amber'
}

const plans: PricingPlan[] = [
  {
    name: 'FREE',
    price: 'Бесплатно',
    period: '',
    description: 'Базовые возможности для начала',
    icon: Zap,
    color: 'cyan',
    buttonText: 'Начать бесплатно',
    features: [
      { text: 'Базовый профиль игрока', included: true },
      { text: 'Статистика текущего сезона', included: true },
      { text: 'Рейтинг по региону', included: true },
      { text: 'Календарь матчей', included: true },
      { text: 'Базовые достижения', included: true },
      { text: 'История за все сезоны', included: false },
      { text: 'Видимость для скаутов', included: false },
      { text: 'AI рекомендации', included: false },
    ],
  },
  {
    name: 'PRO',
    price: '990 ₽',
    period: '/мес',
    description: 'Для серьёзного развития карьеры',
    icon: Sparkles,
    color: 'purple',
    buttonText: 'Выбрать PRO',
    popular: true,
    features: [
      { text: 'Всё из FREE', included: true },
      { text: 'История за все сезоны', included: true },
      { text: 'Графики прогресса', included: true },
      { text: 'Сравнение с игроками', included: true },
      { text: 'Видимость для скаутов', included: true },
      { text: 'Уведомления о просмотрах', included: true },
      { text: 'Фото в профиле', included: true },
      { text: 'AI рекомендации', included: false },
    ],
  },
  {
    name: 'ULTRA',
    price: '2 490 ₽',
    period: '/мес',
    description: 'Максимум возможностей',
    icon: Crown,
    color: 'amber',
    buttonText: 'Выбрать ULTRA',
    features: [
      { text: 'Всё из PRO', included: true },
      { text: 'Приоритет в поиске', included: true },
      { text: 'Персональный URL', included: true },
      { text: 'Видео хайлайты', included: true },
      { text: 'AI рекомендации', included: true },
      { text: 'Экспорт в PDF', included: true },
      { text: 'Сообщения от скаутов', included: true },
      { text: 'Verified badge', included: true },
    ],
  },
]

const colorStyles = {
  cyan: {
    gradient: 'from-[#00d4ff] to-[#0099cc]',
    border: 'border-[#00d4ff]/20',
    glow: 'shadow-[0_0_40px_rgba(0,212,255,0.2)]',
    buttonBg: 'bg-[#00d4ff]/10 hover:bg-[#00d4ff]/20 text-[#00d4ff]',
  },
  purple: {
    gradient: 'from-[#8b5cf6] to-[#6d28d9]',
    border: 'border-[#8b5cf6]/30',
    glow: 'shadow-[0_0_60px_rgba(139,92,246,0.3)]',
    buttonBg: 'bg-gradient-to-r from-[#8b5cf6] to-[#6d28d9] text-white hover:opacity-90',
  },
  amber: {
    gradient: 'from-[#f59e0b] to-[#d97706]',
    border: 'border-[#f59e0b]/20',
    glow: 'shadow-[0_0_40px_rgba(245,158,11,0.2)]',
    buttonBg: 'bg-gradient-to-r from-[#f59e0b] to-[#d97706] text-white hover:opacity-90',
  },
}

const containerVariants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: {
      staggerChildren: 0.15,
    },
  },
}

const itemVariants = {
  hidden: { opacity: 0, y: 40 },
  visible: {
    opacity: 1,
    y: 0,
    transition: {
      type: 'spring',
      damping: 20,
      stiffness: 100,
    },
  },
}

export const LandingPricing = memo(function LandingPricing() {
  return (
    <section
      id="pricing"
      className="relative py-24 px-4 overflow-hidden"
    >
      {/* Background */}
      <div className="absolute inset-0 bg-[#050810]" />

      {/* Decorative elements */}
      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-[600px] h-[600px] bg-[#8b5cf6]/5 rounded-full blur-[150px]" />

      <div className="relative max-w-6xl mx-auto">
        {/* Header */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          whileInView={{ opacity: 1, y: 0 }}
          viewport={{ once: true, margin: '-100px' }}
          transition={{ duration: 0.6 }}
          className="text-center mb-16"
        >
          <h2 className="text-4xl md:text-5xl font-bold text-white mb-4">
            Выбери свой{' '}
            <span className="bg-gradient-to-r from-[#8b5cf6] to-[#00d4ff] bg-clip-text text-transparent">
              тариф
            </span>
          </h2>
          <p className="text-gray-400 text-lg max-w-2xl mx-auto">
            Начни бесплатно и переходи на PRO когда будешь готов
          </p>
        </motion.div>

        {/* Pricing cards */}
        <motion.div
          variants={containerVariants}
          initial="hidden"
          whileInView="visible"
          viewport={{ once: true, margin: '-50px' }}
          className="grid grid-cols-1 md:grid-cols-3 gap-6 lg:gap-8 items-start"
        >
          {plans.map((plan) => {
            const styles = colorStyles[plan.color]
            const Icon = plan.icon

            return (
              <motion.div
                key={plan.name}
                variants={itemVariants}
                className={cn(
                  'relative rounded-2xl p-6 lg:p-8',
                  'bg-white/[0.02] backdrop-blur-sm',
                  'border transition-all duration-300',
                  styles.border,
                  plan.popular && [styles.glow, 'scale-105 lg:scale-110 z-10']
                )}
              >
                {/* Popular badge */}
                {plan.popular && (
                  <div className="absolute -top-3 left-1/2 -translate-x-1/2">
                    <span className={cn(
                      'px-4 py-1 rounded-full text-xs font-semibold',
                      'bg-gradient-to-r',
                      styles.gradient,
                      'text-white'
                    )}>
                      Популярный
                    </span>
                  </div>
                )}

                {/* Header */}
                <div className="text-center mb-6">
                  <div
                    className={cn(
                      'w-14 h-14 rounded-xl flex items-center justify-center mx-auto mb-4',
                      'bg-gradient-to-br',
                      styles.gradient
                    )}
                  >
                    <Icon className="w-7 h-7 text-white" />
                  </div>

                  <h3 className="text-2xl font-bold text-white mb-1">
                    {plan.name}
                  </h3>
                  <p className="text-gray-500 text-sm mb-4">
                    {plan.description}
                  </p>

                  <div className="flex items-baseline justify-center gap-1">
                    <span className="text-4xl font-bold text-white">
                      {plan.price}
                    </span>
                    {plan.period && (
                      <span className="text-gray-500">{plan.period}</span>
                    )}
                  </div>
                </div>

                {/* Features */}
                <ul className="space-y-3 mb-8">
                  {plan.features.map((feature) => (
                    <li
                      key={feature.text}
                      className={cn(
                        'flex items-center gap-3 text-sm',
                        feature.included ? 'text-gray-300' : 'text-gray-600'
                      )}
                    >
                      <Check
                        className={cn(
                          'w-4 h-4 flex-shrink-0',
                          feature.included ? 'text-[#00d4ff]' : 'text-gray-700'
                        )}
                      />
                      {feature.text}
                    </li>
                  ))}
                </ul>

                {/* Button */}
                <button
                  className={cn(
                    'w-full py-3 px-6 rounded-xl font-semibold',
                    'transition-all duration-200',
                    styles.buttonBg
                  )}
                >
                  {plan.buttonText}
                </button>
              </motion.div>
            )
          })}
        </motion.div>

        {/* Bottom note */}
        <motion.p
          initial={{ opacity: 0 }}
          whileInView={{ opacity: 1 }}
          viewport={{ once: true }}
          transition={{ delay: 0.5 }}
          className="text-center text-gray-500 text-sm mt-12"
        >
          Все тарифы можно отменить в любой момент. Оплата через SberPay, Тинькофф или карту.
        </motion.p>
      </div>
    </section>
  )
})
