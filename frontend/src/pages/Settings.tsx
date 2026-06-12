import { useEffect, useMemo, useState } from 'react';
import { Card } from '../components/common/Card';
import { useAppStore } from '../store/useAppStore';
import * as BotService from '../../bindings/github.com/dragonbaba/MyQQBot/internal/service/botservice.js';
import * as ConfigService from '../../bindings/github.com/dragonbaba/MyQQBot/internal/service/configservice.js';
import type { Config, LLMModelInfo, TestLLMResult, ModelCapability } from '../types';
import { Play, Square, TestTube, AlertTriangle, Eye, RefreshCw } from 'lucide-react';

const defaultConfig: Config = {
  llm: {
    base_url: 'https://api.openai.com/v1',
    api_key: '',
    model: 'gpt-4o',
    vision_model: '',
    reasoning_effort: '',
    system_prompt: '',
    model_capabilities: {},
  },
  bot: {
    onebot_ws_url: 'ws://127.0.0.1:3001',
    admin_qq: '',
  },
  search: {
    provider: 'duckduckgo',
    tavily_api_key: '',
  },
};

const reasoningOptions = ['', 'low', 'medium', 'high'];

export function Settings() {
  const setConfig = useAppStore((s) => s.setConfig);
  const addToast = useAppStore((s) => s.addToast);
  const botRunning = useAppStore((s) => s.botStatus.running);
  const [cfg, setLocalConfig] = useState<Config>(defaultConfig);
  const [activeTab, setActiveTab] = useState<'llm' | 'bot' | 'search'>('llm');
  const [starting, setStarting] = useState(false);
  const [testing, setTesting] = useState(false);
  const [testResult, setTestResult] = useState<TestLLMResult | null>(null);
  const [availableModels, setAvailableModels] = useState<LLMModelInfo[]>([]);

  useEffect(() => {
    ConfigService.GetConfig().then((c) => {
      const capabilities: Record<string, ModelCapability> = {};
      if (c.llm.model_capabilities) {
        Object.entries(c.llm.model_capabilities).forEach(([k, v]) => {
          if (v) {
            capabilities[k] = {
              vision: !!v.vision,
              tools: !!v.tools,
              reasoning: !!v.reasoning,
              max_context: v.max_context || 0,
            };
          }
        });
      }
      const loaded: Config = {
        llm: {
          ...defaultConfig.llm,
          ...(c.llm as Config['llm']),
          model_capabilities: capabilities,
        },
        bot: { ...defaultConfig.bot, ...c.bot },
        search: { ...defaultConfig.search, ...c.search },
      };
      setLocalConfig(loaded);
      setConfig(loaded);
    });
  }, [setConfig]);

  const selectedModelInfo = useMemo(() => {
    return availableModels.find((m) => m.id === cfg.llm.model);
  }, [availableModels, cfg.llm.model]);

  const handleSave = async () => {
    try {
      await ConfigService.SaveConfig(cfg as unknown as Parameters<typeof ConfigService.SaveConfig>[0]);
      setConfig(cfg);
      addToast({ type: 'success', message: '配置已保存' });
    } catch (err) {
      addToast({ type: 'error', message: `保存失败: ${err}` });
    }
  };

  const handleTestLLM = async () => {
    setTesting(true);
    setTestResult(null);
    try {
      const result = await ConfigService.TestLLMConnection(cfg.llm.base_url, cfg.llm.api_key);
      setTestResult(result);
      if (result.success) {
        setAvailableModels(result.models);
        addToast({ type: 'success', message: `连接正常，发现 ${result.models.length} 个模型` });
      } else {
        addToast({ type: 'error', message: result.error });
      }
    } catch (err) {
      addToast({ type: 'error', message: `连接失败: ${err}` });
    } finally {
      setTesting(false);
    }
  };

  const handleToggleBot = async () => {
    setStarting(true);
    try {
      if (botRunning) {
        await BotService.StopBot();
        addToast({ type: 'success', message: '机器人已停止' });
      } else {
        await BotService.StartBot();
        addToast({ type: 'success', message: '机器人已启动' });
      }
      const newStatus = await BotService.GetBotStatus();
      useAppStore.getState().setBotStatus(newStatus);
    } catch (err) {
      addToast({ type: 'error', message: `操作失败: ${err}` });
    } finally {
      setStarting(false);
    }
  };

  const updateLLM = (patch: Partial<Config['llm']>) =>
    setLocalConfig((prev) => ({ ...prev, llm: { ...prev.llm, ...patch } }));
  const updateBot = (patch: Partial<Config['bot']>) =>
    setLocalConfig((prev) => ({ ...prev, bot: { ...prev.bot, ...patch } }));
  const updateSearch = (patch: Partial<Config['search']>) =>
    setLocalConfig((prev) => ({ ...prev, search: { ...prev.search, ...patch } }));

  const needsVisionFallback = selectedModelInfo && !selectedModelInfo.vision;

  return (
    <div className="h-full overflow-auto p-6">
      <Card className="max-w-3xl mx-auto">
        <div className="flex gap-2 mb-6 border-b border-border-DEFAULT pb-2">
          {(['llm', 'bot', 'search'] as const).map((tab) => (
            <button
              key={tab}
              onClick={() => setActiveTab(tab)}
              className={`px-4 py-2 text-sm font-medium rounded-lg transition-colors cursor-pointer ${
                activeTab === tab
                  ? 'bg-brand-cta/15 text-brand-cta'
                  : 'text-text-secondary hover:text-text-primary hover:bg-surface-elevated'
              }`}
            >
              {tab === 'llm' && 'LLM 配置'}
              {tab === 'bot' && 'QQ 机器人'}
              {tab === 'search' && '搜索设置'}
            </button>
          ))}
        </div>

        <div className="space-y-4">
          {activeTab === 'llm' && (
            <>
              <label className="block text-sm text-text-secondary">
                Base URL
                <input
                  value={cfg.llm.base_url}
                  onChange={(e) => updateLLM({ base_url: e.target.value })}
                  placeholder="https://api.openai.com/v1"
                  className="mt-1 w-full bg-surface-base border border-border-subtle rounded-lg px-3 py-2 text-text-primary focus:outline-none focus:border-brand-cta/50"
                />
              </label>
              <label className="block text-sm text-text-secondary">
                API Key
                <input
                  type="password"
                  value={cfg.llm.api_key}
                  onChange={(e) => updateLLM({ api_key: e.target.value })}
                  className="mt-1 w-full bg-surface-base border border-border-subtle rounded-lg px-3 py-2 text-text-primary focus:outline-none focus:border-brand-cta/50"
                />
              </label>

              <div className="flex items-end gap-2">
                <label className="block text-sm text-text-secondary flex-1">
                  Model
                  <select
                    value={cfg.llm.model}
                    onChange={(e) => updateLLM({ model: e.target.value })}
                    className="mt-1 w-full bg-surface-base border border-border-subtle rounded-lg px-3 py-2 text-text-primary focus:outline-none focus:border-brand-cta/50"
                  >
                    <option value="">请选择或手动输入</option>
                    {availableModels.map((m) => (
                      <option key={m.id} value={m.id}>
                        {m.id}
                        {m.vision ? ' · 视觉' : ''}
                        {m.tools ? ' · 工具' : ''}
                        {m.reasoning ? ' · 推理' : ''}
                      </option>
                    ))}
                  </select>
                </label>
                <input
                  value={cfg.llm.model}
                  onChange={(e) => updateLLM({ model: e.target.value })}
                  placeholder="手动输入模型名"
                  className="w-64 bg-surface-base border border-border-subtle rounded-lg px-3 py-2 text-sm text-text-primary focus:outline-none focus:border-brand-cta/50"
                />
              </div>

              {selectedModelInfo && (
                <div className="flex flex-wrap gap-2 text-xs">
                  <span className={`px-2 py-1 rounded-md ${selectedModelInfo.vision ? 'bg-status-online/10 text-status-online' : 'bg-status-error/10 text-status-error'}`}>
                    视觉: {selectedModelInfo.vision ? '支持' : '不支持'}
                  </span>
                  <span className={`px-2 py-1 rounded-md ${selectedModelInfo.tools ? 'bg-status-online/10 text-status-online' : 'bg-status-error/10 text-status-error'}`}>
                    工具: {selectedModelInfo.tools ? '支持' : '不支持'}
                  </span>
                  <span className={`px-2 py-1 rounded-md ${selectedModelInfo.reasoning ? 'bg-status-online/10 text-status-online' : 'bg-status-warning/10 text-status-warning'}`}>
                    推理: {selectedModelInfo.reasoning ? '支持' : '不支持'}
                  </span>
                </div>
              )}

              {needsVisionFallback && (
                <div className="p-3 rounded-lg bg-status-warning/10 border border-status-warning/20 flex items-start gap-2">
                  <AlertTriangle className="w-4 h-4 text-status-warning shrink-0 mt-0.5" />
                  <p className="text-xs text-status-warning">
                    当前模型不支持视觉处理。如需处理图片消息，请在下方的「视觉辅助模型」中指定一个支持视觉的模型（如 gpt-4o）。
                  </p>
                </div>
              )}

              <label className="block text-sm text-text-secondary">
                <span className="flex items-center gap-1">
                  <Eye className="w-4 h-4" /> 视觉辅助模型
                </span>
                <input
                  value={cfg.llm.vision_model}
                  onChange={(e) => updateLLM({ vision_model: e.target.value })}
                  placeholder={needsVisionFallback ? '必填：如 gpt-4o' : '可选'}
                  className="mt-1 w-full bg-surface-base border border-border-subtle rounded-lg px-3 py-2 text-text-primary focus:outline-none focus:border-brand-cta/50"
                />
              </label>

              <label className="block text-sm text-text-secondary">
                思考强度（仅部分推理模型有效）
                <select
                  value={cfg.llm.reasoning_effort}
                  onChange={(e) => updateLLM({ reasoning_effort: e.target.value })}
                  className="mt-1 w-full bg-surface-base border border-border-subtle rounded-lg px-3 py-2 text-text-primary focus:outline-none focus:border-brand-cta/50"
                >
                  {reasoningOptions.map((o) => (
                    <option key={o} value={o}>
                      {o === '' ? '默认' : o}
                    </option>
                  ))}
                </select>
              </label>

              <label className="block text-sm text-text-secondary">
                System Prompt
                <textarea
                  value={cfg.llm.system_prompt}
                  onChange={(e) => updateLLM({ system_prompt: e.target.value })}
                  rows={4}
                  className="mt-1 w-full bg-surface-base border border-border-subtle rounded-lg px-3 py-2 text-text-primary focus:outline-none focus:border-brand-cta/50 resize-none"
                />
              </label>

              <div className="flex items-center gap-2">
                <button
                  onClick={handleTestLLM}
                  disabled={testing}
                  className="inline-flex items-center gap-2 px-4 py-2 bg-status-info/15 text-status-info hover:bg-status-info/25 rounded-lg transition-colors cursor-pointer disabled:opacity-50"
                >
                  {testing ? <RefreshCw className="w-4 h-4 animate-spin" /> : <TestTube className="w-4 h-4" />}
                  {testing ? '测试中...' : '测试并获取模型列表'}
                </button>
                {testResult && !testResult.success && (
                  <span className="text-xs text-status-error">{testResult.error}</span>
                )}
              </div>
            </>
          )}

          {activeTab === 'bot' && (
            <>
              <label className="block text-sm text-text-secondary">
                OneBot WebSocket URL
                <input
                  value={cfg.bot.onebot_ws_url}
                  onChange={(e) => updateBot({ onebot_ws_url: e.target.value })}
                  className="mt-1 w-full bg-surface-base border border-border-subtle rounded-lg px-3 py-2 text-text-primary focus:outline-none focus:border-brand-cta/50"
                />
              </label>
              <label className="block text-sm text-text-secondary">
                Admin QQ
                <input
                  value={cfg.bot.admin_qq}
                  onChange={(e) => updateBot({ admin_qq: e.target.value })}
                  className="mt-1 w-full bg-surface-base border border-border-subtle rounded-lg px-3 py-2 text-text-primary focus:outline-none focus:border-brand-cta/50"
                />
              </label>
              <button
                onClick={handleToggleBot}
                disabled={starting}
                className="inline-flex items-center gap-2 px-4 py-2 bg-brand-cta/15 text-brand-cta hover:bg-brand-cta/25 rounded-lg transition-colors cursor-pointer disabled:opacity-50"
              >
                {botRunning ? (
                  <>
                    <Square className="w-4 h-4" /> 停止机器人
                  </>
                ) : (
                  <>
                    <Play className="w-4 h-4" /> 启动机器人
                  </>
                )}
              </button>
            </>
          )}

          {activeTab === 'search' && (
            <>
              <div className="space-y-2">
                <p className="text-sm text-text-secondary">搜索 Provider</p>
                <label className="flex items-center gap-2 text-sm text-text-primary cursor-pointer">
                  <input
                    type="radio"
                    name="provider"
                    checked={cfg.search.provider === 'duckduckgo'}
                    onChange={() => updateSearch({ provider: 'duckduckgo' })}
                    className="cursor-pointer"
                  />
                  DuckDuckGo（免费）
                </label>
                <label className="flex items-center gap-2 text-sm text-text-primary cursor-pointer">
                  <input
                    type="radio"
                    name="provider"
                    checked={cfg.search.provider === 'tavily'}
                    onChange={() => updateSearch({ provider: 'tavily' })}
                    className="cursor-pointer"
                  />
                  Tavily
                </label>
              </div>
              {cfg.search.provider === 'tavily' && (
                <label className="block text-sm text-text-secondary">
                  Tavily API Key
                  <input
                    type="password"
                    value={cfg.search.tavily_api_key}
                    onChange={(e) => updateSearch({ tavily_api_key: e.target.value })}
                    className="mt-1 w-full bg-surface-base border border-border-subtle rounded-lg px-3 py-2 text-text-primary focus:outline-none focus:border-brand-cta/50"
                  />
                </label>
              )}
            </>
          )}
        </div>

        <div className="mt-6 flex justify-end">
          <button
            onClick={handleSave}
            className="px-6 py-2 bg-brand-cta hover:bg-brand-cta/80 text-surface-base font-medium rounded-lg transition-colors cursor-pointer"
          >
            保存配置
          </button>
        </div>
      </Card>
    </div>
  );
}
