{{- define "ingressAnnotations" }}
  {{- $annotations := dict }}
  {{- $componentSpec := index .Values .component }}

  {{- if hasKey $componentSpec.ingress "annotations" }}
    {{- $_ := merge $annotations $componentSpec.ingress.annotations }}
  {{- end }}

  {{- $_ := merge $annotations .Values.ingress.annotations .Values.jxRequirements.ingress.annotations  }}

  {{- if not (hasKey $annotations "kubernetes.io/ingress.class") }}
    {{- $customIngressClass := "" }}
    {{- if $componentSpec.ingress.customIngressClass }}
      {{- $customIngressClass := $componentSpec.ingress.customIngressClass }}
      {{- $_ := set $annotations "kubernetes.io/ingress.class" $customIngressClass }}
    {{- else if hasKey .Values.ingress "customIngressClass" }}
      {{- if eq .component "docker-registry" }}
        {{- if hasKey .Values.ingress.customIngressClass "dockerRegistry" }}
          {{- $customIngressClass := index .Values.ingress.customIngressClass "dockerRegistry" }}
          {{- $_ := set $annotations "kubernetes.io/ingress.class" $customIngressClass }}
        {{- end }}
      {{- else if hasKey .Values.ingress.customIngressClass .component }}
        {{- $customIngressClass := index .Values.ingress.customIngressClass .component }}
        {{- $_ := set $annotations "kubernetes.io/ingress.class" $customIngressClass }}
      {{- end }}
    {{- end }}
    {{- if not (hasKey $annotations "kubernetes.io/ingress.class") }}
      {{- $_ := set $annotations "kubernetes.io/ingress.class" ($customIngressClass | default "nginx")  }}
    {{- end }}
  {{- end }}
  {{- if and (hasKey .Values.jxRequirements.ingress "serviceType") (.Values.jxRequirements.ingress.serviceType) (eq .Values.jxRequirements.ingress.serviceType "NodePort") (not (hasKey $annotations "jenkins.io/host")) }}
    {{- $_ := set $annotations "jenkins.io/host" .Values.jxRequirements.ingress.domain }}
  {{- end }}
  {{- if $annotations }}
{{ toYaml $annotations | indent 4 }}
  {{- end }}
{{- end }}

{{- /*
httpRouteParentRefs renders the parentRefs list for an HTTPRoute.

When jxRequirements.ingress.kind is "httproute" the default envoy-gateway "http"
listener is attached, plus the "https" listener when
jxRequirements.ingress.tls.enabled is true. Any per-component
httpRoute.customParentRefs are then appended.

When ingress.kind is anything other than "httproute" no defaults are injected
and the parentRefs are exactly the user-supplied httpRoute.customParentRefs, so
users running their own gateway can fully control attachment.

Call with (dict "Values" .Values "component" "<name>").
*/ -}}
{{- define "httpRouteParentRefs" -}}
{{- $spec := index .Values .component -}}
{{- $refs := list -}}
{{- if eq "httproute" .Values.jxRequirements.ingress.kind -}}
{{- $refs = append $refs (dict "name" "envoy-gateway" "namespace" "envoy-gateway-system" "sectionName" "http") -}}
{{- if .Values.jxRequirements.ingress.tls.enabled -}}
{{- $refs = append $refs (dict "name" "envoy-gateway" "namespace" "envoy-gateway-system" "sectionName" "https") -}}
{{- end -}}
{{- end -}}
{{- range $spec.httpRoute.customParentRefs -}}
{{- $refs = append $refs . -}}
{{- end -}}
{{- if $refs -}}
{{ toYaml $refs }}
{{- end -}}
{{- end }}

{{- /*
httpRouteHostname returns the hostname for a component's HTTPRoute: the
per-component httpRoute.customHost if set, otherwise the composed
<prefix><namespaceSubDomain><domain>. The prefix falls back to the component's
ingress.prefix when httpRoute.prefix is not set. Call with
(dict "Values" .Values "component" "<name>").
*/ -}}
{{- define "httpRouteHostname" -}}
{{- $spec := index .Values .component -}}
{{- if $spec.httpRoute.customHost -}}
{{ $spec.httpRoute.customHost }}
{{- else -}}
{{ $spec.httpRoute.prefix | default $spec.ingress.prefix }}{{ .Values.jxRequirements.ingress.namespaceSubDomain }}{{ .Values.jxRequirements.ingress.domain }}
{{- end -}}
{{- end }}

{{- /*
httpRouteAnnotations returns the per-component httpRoute.annotations as raw
(unindented) YAML, or an empty string when there are none. Call with
(dict "Values" .Values "component" "<name>"); the caller is responsible for
indenting (e.g. `nindent 4`) and for only emitting the `annotations:` key when
the result is non-empty.
*/ -}}
{{- define "httpRouteAnnotations" -}}
{{- $spec := index .Values .component -}}
{{- if $spec.httpRoute.annotations -}}
{{ toYaml $spec.httpRoute.annotations }}
{{- end -}}
{{- end }}
