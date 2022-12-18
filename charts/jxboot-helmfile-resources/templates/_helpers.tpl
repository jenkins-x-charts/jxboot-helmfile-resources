{{- define "ingressAnnotations" }}
{{- $annotations := dict }}
{{- $_ := merge $annotations .Values.ingress.annotations .Values.jxRequirements.ingress.annotations  }}
{{- if not (hasKey $annotations "kubernetes.io/ingress.class") }}
{{- $customIngressClass := index .Values.ingress.customIngressClass .component }}
{{- $_ := set $annotations "kubernetes.io/ingress.class" ($customIngressClass | default "nginx") }}
{{- end }}
{{- if and (hasKey .Values.jxRequirements.ingress "serviceType") (.Values.jxRequirements.ingress.serviceType) (eq .Values.jxRequirements.ingress.serviceType "NodePort") (not (hasKey $annotations "jenkins.io/host")) }}
{{- $_ := set $annotations "jenkins.io/host" .Values.jxRequirements.ingress.domain }}
{{- end }}
{{- $customIngressAnnotation := index .Values.ingress.customIngressAnnotation .component }}
{{- if $annotations }}
{{- $_ := merge $annotations $customIngressAnnotation }}
{{ toYaml $annotations | indent 4 }}
{{- end }}
{{- end -}}